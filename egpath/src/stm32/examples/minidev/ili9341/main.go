package main

import (
	"delay"
	"fmt"
	"rtos"

	"display/ili9341"
	"display/ili9341/ili9341test"

	"stm32/ilidci"

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
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7

	gpio.B.EnableClock(true)
	ilidc := gpio.B.Pin(0)
	ilireset := gpio.B.Pin(1)
	ilics := gpio.B.Pin(10)

	// SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	d := dma.DMA1
	d.EnableClock(true)
	lcdspi = spi.NewDriver(spi.SPI1, d.Channel(3, 0), d.Channel(2, 0))
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	// Controll

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

	lcd = ili9341.NewDisplay(ilidci.New(lcdspi, ilidc, 36e6), 240, 320)
	lcd.DCI().Setup()
}

func main() {
	delay.Millisec(100) // For SWO output.

	spip := lcdspi.Periph()
	fmt.Printf("\nSPI on %s (%d MHz).\n", spip.Bus(), spip.Bus().Clock()/1e6)
	fmt.Printf("SPI speed: %d bps.\n", spip.Baudrate(spip.Conf()))
	ili9341test.Run(lcd, 10, true)
}

func lcdSPIISR() {
	lcdspi.ISR()
}

func lcdRxDMAISR() {
	lcdspi.DMAISR(lcdspi.RxDMA())
}

func lcdTxDMAISR() {
	lcdspi.DMAISR(lcdspi.TxDMA())
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,

	irq.SPI1:          lcdSPIISR,
	irq.DMA1_Channel2: lcdRxDMAISR,
	irq.DMA1_Channel3: lcdTxDMAISR,
}
