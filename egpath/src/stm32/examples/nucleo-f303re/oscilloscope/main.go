package main

import (
	"delay"
	"image"
	"rtos"

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
)

var (
	lcdspi *spi.Driver
	lcd    *ili9341.Display
	inp1   *adc.ADC_Periph
)

const advregen = 1 << adc.ADVREGENn

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

	// ILI SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	d := dma.DMA1
	d.EnableClock(true)
	lcdspi = spi.NewDriver(spi.SPI1, d.Channel(2, 0), d.Channel(3, 0))
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
	// Select AHBclk/div as ADC12 clock source.
	const log2div = 0 // Max. 2.
	adc.ADC1_2.CKMODE().Store((log2div + 1) << adc.CKMODEn)
	inp1 = adc.ADC1
	// Enable voltage regulator.
	inp1.CR.Store(0)
	inp1.CR.Store(advregen)
	delay.Millisec(1)
	// Start calibration.
	inp1.CR.Store(adc.ADCAL | advregen)
	for inp1.ADCAL().Load() != 0 {
	}
	// ADEN can be set not sooner than 4 ADC clock cycles after ADCAL == 0.
	delay.Loop(5 << log2div)
	inp1.CR.Store(adc.ADEN | advregen)
	// 8-bit resolution.
	inp1.RES().Store(2 << adc.RESn)
	// Use only CH1.
	inp1.SQR1.Store(1<<adc.SQ1n | (1-1)<<adc.Ln)
}

func load(inp *adc.ADC_Periph) int {
	inp.CR.Store(adc.ADSTART | advregen)
	for inp.EOC().Load() == 0 {
	}
	return int(inp1.DR.Load())
}

func main() {
	lcd.SlpOut()
	delay.Millisec(120)
	lcd.DispOn()
	lcd.PixSet(ili9341.PF16) // 16-bit pixel format.
	lcd.MADCtl(ili9341.MY | ili9341.MX | ili9341.MV | ili9341.BGR)
	lcd.SetWordSize(16)

	lcd.SetColor(0)
	lcd.Rect(lcd.Bounds())

	// Wait for ADC.
	for inp1.ADRDY().Load() == 0 {
	}

	wh := lcd.Bounds().Max
	buf := make([]byte, wh.X)
	trig := 128
	for {
		for load(inp1) >= trig {
		}
		var first int
		for {
			first = load(inp1)
			if first > trig {
				break
			}
		}
		buf[0] = byte(first)
		for x := 1; x < wh.X; x++ {
			buf[x] = byte(load(inp1))
		}
		lcd.SetColor(0)
		lcd.Rect(image.Rect(0, 0, 1, wh.Y))
		for x := 1; x < wh.X; x++ {
			y0 := wh.Y - 8 - int(buf[x-1])*7/8
			y1 := wh.Y - 8 - int(buf[x])*7/8
			lcd.SetColor(0)
			lcd.Rect(image.Rect(x, 0, x+1, wh.Y))
			lcd.SetColor(0xffff)
			lcd.Line(image.Pt(x-1, y0), image.Pt(x, y1))
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

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI1:          lcdSPIISR,
	irq.DMA1_Channel2: lcdRxDMAISR,
	irq.DMA1_Channel3: lcdTxDMAISR,
}
