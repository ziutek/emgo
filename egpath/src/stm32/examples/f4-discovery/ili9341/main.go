// This is test program that tries to communicate with ILI9341 controller using
// raw SPI (without any display library). There are more useful examples (for
// other boards) that use display/ili9341 package.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

type ILI9341 struct {
	SPI *spi.Driver

	ctrport       *gpio.Port
	reset, cs, dc gpio.Pins
}

func (ili *ILI9341) Reset() {
	ili.ctrport.SetPins(ili.reset | ili.cs)
	delay.Millisec(1)
	ili.ctrport.ClearPins(ili.reset)
	delay.Millisec(1)
	ili.ctrport.SetPins(ili.reset)
}

func (ili *ILI9341) Select() {
	ili.ctrport.ClearPins(ili.cs)
}

func (ili *ILI9341) Deselect() {
	ili.ctrport.SetPins(ili.cs)
}

func (ili *ILI9341) Cmd(cmd byte) {
	delay.Millisec(1)
	ili.ctrport.ClearPins(ili.dc)
	ili.SPI.WriteReadByte(cmd)
	delay.Millisec(1)
	ili.ctrport.SetPins(ili.dc)
}

func (ili *ILI9341) WriteByte(b byte) {
	ili.SPI.WriteReadByte(b)
}

func (ili *ILI9341) Write(data []byte) {
	ili.SPI.WriteRead(data, nil)
}

var ili ILI9341

func init() {
	system.Setup168(8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	ili.ctrport = gpio.A
	ili.reset, ili.cs, ili.dc = gpio.Pin0, gpio.Pin1, gpio.Pin4
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7

	// SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	d := dma.DMA2
	d.EnableClock(true)
	ili.SPI = spi.NewDriver(spi.SPI1, d.Channel(2, 3), d.Channel(3, 3))
	ili.SPI.P.EnableClock(true)
	ili.SPI.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			ili.SPI.P.BR(21e6) | // 21 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	ili.SPI.P.Enable()
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA2_Stream2).Enable()
	rtos.IRQ(irq.DMA2_Stream3).Enable()

	// Controll

	ili.ctrport.Setup(
		ili.reset|ili.cs|ili.dc,
		&gpio.Config{Mode: gpio.Out, Speed: gpio.High},
	)
}

const (
	ILI9341_SLPOUT = 0x11
	ILI9341_DISPON = 0x29
	ILI9341_RAMWR  = 0x2C
	ILI9341_MADCTL = 0x36
	ILI9341_PIXFMT = 0x3A
)

func main() {
	ili.Reset()
	delay.Millisec(10)
	ili.Select()
	ili.Cmd(ILI9341_SLPOUT)
	delay.Millisec(120)
	ili.Cmd(ILI9341_DISPON)

	spibus := ili.SPI.P.Bus()
	fmt.Printf(
		"\nSPI on %s (%d MHz). SPI speed: %d Hz.\n\n",
		spibus, spibus.Clock()/1e6, ili.SPI.P.Baudrate(ili.SPI.P.Conf()),
	)

	ili.Cmd(ILI9341_PIXFMT)
	ili.WriteByte(0x55) // 16 bit 565 format.

	ili.Cmd(ILI9341_MADCTL)
	ili.WriteByte(0x48) // Screen orientation.

	ili.Cmd(ILI9341_RAMWR)

	line := make([]byte, 240*2)
	for i := 0; i < 320; i++ {
		for k := range line {
			x := i
			if k == x || k == x+1 {
				line[k] = 0xff
			} else {
				line[k] = 0
			}
		}
		ili.Write(line)
	}
}

func iliSPIISR() {
	ili.SPI.ISR()
}

func iliRxDMAISR() {
	ili.SPI.DMAISR(ili.SPI.RxDMA)
}

func iliTxDMAISR() {
	ili.SPI.DMAISR(ili.SPI.TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI1:         iliSPIISR,
	irq.DMA2_Stream2: iliRxDMAISR,
	irq.DMA2_Stream3: iliTxDMAISR,
}
