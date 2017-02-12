package main

import (
	"delay"
	"image"
	"rtos"

	"display/ili9341"

	"stm32/ilidci"

	"stm32/hal/adc"
	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/opamp"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

var (
	lcdspi *spi.Driver
	lcd    *ili9341.Display
	adcd   *adc.Driver
)

func checkErr(err error) {
	if err == nil {
		return
	}
	rtos.Debug(0).WriteString(err.Error())
	for {
	}
}

func init() {
	system.Setup(8, 1, 72/8)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7
	adcin := gpio.A.Pin(0)   // ADC1 input.
	opampin := gpio.A.Pin(1) // OPAMP1 input.
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

	adcin.Setup(&gpio.Config{Mode: gpio.Ana})

	rcc.RCC.TIM6EN().Set()
	adcd = adc.NewDriver(adc.ADC1, dma1.Channel(1, 0))
	adcd.P.EnableClock(true)
	adcd.P.EnableVoltage()
	delay.Millisec(1)
	adcd.P.SetClockMode(adc.HCLK1)
	adcd.P.Callibrate()
	// delay.Loop(5 << (adc.HCLK1 - 1)) // Need before adcd.Enable()

	rtos.IRQ(irq.ADC1_2).Enable()
	rtos.IRQ(irq.DMA1_Channel1).Enable()

	// ADC operational amplifier

	opampin.Setup(&gpio.Config{Mode: gpio.Ana})

	rcc.RCC.SYSCFGEN().Set()
	opamp.OPAMP1.CSR.Store(opamp.OPAMPxEN |
		3<<opamp.VPSELn | // Positive input connected to PA1.
		3<<opamp.VMSELn | // 2: PGA mode, 3: follower mode.
		0<<opamp.PGGAINn, // Gain: 0:2, 1:4, 2:8, 3:16.
	)

	// ADC timer.

	div1, div2 := 2, 5 // ADC SR = 72 MHz / div1 / div2
	t := tim.TIM6
	t.CR2.Store(2 << tim.MMSn)
	t.ARR.Store(tim.ARR_Bits(div2 - 1))
	t.PSC.Store(tim.PSC_Bits(div1 - 1))
	t.CR1.Store(tim.ARPE | tim.CEN)
}

func main() {
	lcd.SlpOut()
	delay.Millisec(120)
	lcd.DispOn()
	lcd.PixSet(ili9341.PF16) // 16-bit pixel format.
	lcd.MADCtl(ili9341.MY | ili9341.MX | ili9341.MV | ili9341.BGR)
	lcd.SetWordSize(16)

	scr := lcd.Area(lcd.Bounds())

	scr.SetColor(0)
	scr.FillRect(scr.Bounds())

	adcd.P.SetResolution(adc.Res8)
	//adcd.P.SetRegularSeq(1) // PA0
	adcd.P.SetRegularSeq(3) // PA2 (OPAMP1)
	adcd.P.SetExtTrigSrc(adc.ADC12_TIM6_TRGO)
	adcd.P.SetExtTrigEdge(adc.EdgeRising)
	checkErr(adcd.Enable())

	wh := scr.Bounds().Max
	scale := func(y byte) int { return wh.Y - 8 - int(y)*7/8 }
	buf := make([]byte, wh.X*3)
	const trig = 128
	for {
		_, err := adcd.Read(buf)
		checkErr(err)
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
			scr.SetColor(0)
			scr.FillRect(image.Rect(x, 0, x+1, wh.Y))
			scr.SetColor(0xffff)
			y0 := scale(buf[offset+x])
			y1 := scale(buf[offset+x+1])
			if y0 > y1 {
				y0, y1 = y1, y0
			}
			y1++
			scr.FillRect(image.Rectangle{image.Pt(x, y0), image.Pt(x+1, y1)})
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

func adcISR() {
	adcd.ISR()
}

func adcDMAISR() {
	adcd.DMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI1:          lcdSPIISR,
	irq.DMA1_Channel2: lcdRxDMAISR,
	irq.DMA1_Channel3: lcdTxDMAISR,

	irq.ADC1_2:        adcISR,
	irq.DMA1_Channel1: adcDMAISR,
}
