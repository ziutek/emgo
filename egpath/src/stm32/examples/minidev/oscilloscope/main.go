package main

import (
	"delay"
	//"image"
	//"image/color"
	"rtos"

	"display/ili9341"

	"stm32/ilidci"

	"stm32/hal/adc"
	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
)

var (
	lcdspi *spi.Driver
	lcd    *ili9341.Display
	adcd   *adc.Driver
)

func init() {
	system.Setup(8, 1, 72/8)
	rtcst.Setup(32768)

	// GPIO

	gpio.A.EnableClock(true)
	adcin := gpio.A.Pin(0)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7

	gpio.B.EnableClock(true)
	ilidc := gpio.B.Pin(0)
	ilireset := gpio.B.Pin(1)
	ilics := gpio.B.Pin(10)

	// DMA
	dma1 := dma.DMA1
	dma1.EnableClock(true)

	// ILI SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
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

	// ILI controll

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
	adcd = adc.NewDriver(adc.ADC1, dma1.Channel(1, 0))
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
	irq.RTCAlarm: rtcst.ISR,

	irq.SPI1:          lcdSPIISR,
	irq.DMA1_Channel2: lcdRxDMAISR,
	irq.DMA1_Channel3: lcdTxDMAISR,
}
