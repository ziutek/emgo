package main

import (
	"delay"
	"image"
	"rtos"

	//"arch/cortexm/debug/itm"

	"display/ili9341"

	"stm32/ilidci"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/adc"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

var (
	lcdspi *spi.Driver
	lcd    *ili9341.Display
	adcdrv *ADCDriver
)

func init() {
	system.Setup(8, 1, 72/8)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7
	ain := gpio.A.Pin(0)
	ilics := gpio.A.Pin(15)

	gpio.B.EnableClock(true)
	ilidc := gpio.B.Pin(7)

	gpio.C.EnableClock(true)
	//spiport, sck, miso, mosi := gpio.C, gpio.Pin10, gpio.Pin11, gpio.Pin12
	ilireset := gpio.C.Pin(13) // Max output: 2 MHz, 3 mA.

	// DMA
	dma1 := dma.DMA1
	dma1.EnableClock(true)

	// ILI SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	lcdspi = spi.NewDriver(spi.SPI1, dma1.Channel(2, 0), dma1.Channel(3, 0))
	lcdspi.P.EnableClock(true)
	lcdspi.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			lcdspi.P.BR(36e6) | // 36 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	lcdspi.P.SetWordSize(8)
	lcdspi.P.Enable()
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	// ILI Controll

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.High}
	ilics.Setup(&cfg)
	ilics.Set()
	ilidc.Setup(&cfg)
	cfg.Speed = gpio.Low
	ilireset.Setup(&cfg)
	delay.Millisec(1) // Reset pulse.
	ilireset.Set()
	delay.Millisec(5) // Wait for reset.
	ilics.Clear()

	lcd = ili9341.NewDisplay(ilidci.NewDCI(lcdspi, ilidc))

	// ADC

	ain.Setup(&gpio.Config{Mode: gpio.Ana})

	rcc.RCC.ADC12EN().Set()
	rcc.RCC.TIM6EN().Set()

	adcdrv = NewADCDriver(adc.ADC1, dma1.Channel(1, 0), tim.TIM6)
	adcdrv.EnableVReg()
	adcdrv.Callibrate(0)
	adcdrv.Enable()
	adcdrv.SetResolution(Res8)
	adcdrv.SetRegularSeq(1)
	adcdrv.SetExtTrigSrc(ADC12_TIM6_TRGO)
	adcdrv.SetExtTrigEdge(EdgeRising)
	adcdrv.SetTimer(2, 5) // 72 MHz / (2 * 5) = 7.2 MHz (max. for 8 bit res.)
	adcdrv.StartADC()

	rtos.IRQ(irq.DMA1_Channel1).Enable()
}

func main() {
	lcd.SlpOut()
	delay.Millisec(120)
	lcd.DispOn()
	lcd.PixSet(ili9341.PF16) // 16-bit pixel format.
	lcd.MADCtl(ili9341.MY | ili9341.MX | ili9341.MV | ili9341.BGR)
	lcd.SetWordSize(16)

	lcd.SetColor(0)
	lcd.FillRect(lcd.Bounds())

	wh := lcd.Bounds().Max
	scale := func(y byte) int { return wh.Y - 8 - int(y)*7/8 }
	buf := make([]byte, wh.X*3)
	const trig = 128
	for {
		adcdrv.Read(buf)
		offset := -1
		for i, b := range buf[:wh.X*2] {
			if b < trig {
				if buf[i+1] >= trig {
					offset = i
					break
				}
			}
		}
		if offset < 0 {
			offset = 0
		}
		for x := 0; x < wh.X; x++ {
			lcd.SetColor(0)
			lcd.FillRect(image.Rect(x, 0, x+1, wh.Y))
			lcd.SetColor(0xffff)
			y0 := scale(buf[offset+x])
			y1 := scale(buf[offset+x+1])
			if y0 > y1 {
				y0, y1 = y1, y0
			}
			y1++
			lcd.FillRect(image.Rectangle{image.Pt(x, y0), image.Pt(x+1, y1)})
		}
	}
}

func lcdSPIISR() {
	lcdspi.ISR()
}

func lcdRxDMAISR() {
	lcdspi.DMAISR(lcdspi.RxDMA)
}

func lcdTxDMAISR() {
	lcdspi.DMAISR(lcdspi.TxDMA)
}

func adcDMAISR() {
	adcdrv.DMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI1:          lcdSPIISR,
	irq.DMA1_Channel2: lcdRxDMAISR,
	irq.DMA1_Channel3: lcdTxDMAISR,

	irq.DMA1_Channel1: adcDMAISR,
}
