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

	reset, cs, dc gpio.Pin
}

func (ili *ILI9341) Reset() {
	ili.reset.Set()
	ili.cs.Set()
	delay.Millisec(1)
	ili.reset.Clear()
	delay.Millisec(1)
	ili.reset.Set()
}

func (ili *ILI9341) Select() {
	ili.cs.Clear()
}

func (ili *ILI9341) Deselect() {
	ili.cs.Set()
}

func (ili *ILI9341) Cmd(cmd byte) {
	delay.Millisec(1)
	ili.dc.Clear()
	ili.SPI.WriteReadByte(cmd)
	delay.Millisec(1)
	ili.dc.Set()
}

func (ili *ILI9341) WriteByte(b byte) {
	ili.SPI.WriteReadByte(b)
}

func (ili *ILI9341) Write(data []byte) {
	ili.SPI.WriteRead(data, nil)
}

var ili ILI9341

func init() {
	system.Setup(8, 1, 72/8)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	ili.dc = gpio.A.Pin(15)

	gpio.B.EnableClock(true)
	ili.cs = gpio.B.Pin(7)

	gpio.C.EnableClock(true)
	spiport, sck, miso, mosi := gpio.C, gpio.Pin10, gpio.Pin11, gpio.Pin12
	ili.reset = gpio.C.Pin(13) // Max output: 2 MHz, 3 mA.

	// SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI3_AF6)
	d := dma.DMA2
	d.EnableClock(true)
	ili.SPI = spi.NewDriver(spi.SPI3, d.Channel(1, 0), d.Channel(2, 0))
	ili.SPI.P.EnableClock(true)
	ili.SPI.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			ili.SPI.P.BR(18e6) | // 18 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	ili.SPI.P.SetWordSize(8)
	ili.SPI.P.Enable()
	rtos.IRQ(irq.SPI3).Enable()
	rtos.IRQ(irq.DMA2_Channel1).Enable()
	rtos.IRQ(irq.DMA2_Channel2).Enable()

	// Controll

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.High}
	ili.cs.Setup(&cfg)
	ili.dc.Setup(&cfg)
	cfg.Speed = gpio.Low
	ili.reset.Setup(&cfg)
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

	spibus := ili.SPI.P.Bus()
	fmt.Printf(
		"\nSPI on %s (%d MHz). SPI speed: %d Hz.\n\n",
		spibus, spibus.Clock()/1e6, ili.SPI.P.Baudrate(ili.SPI.P.Conf()),
	)

	ili.Select()
	ili.Cmd(ILI9341_SLPOUT)
	delay.Millisec(120)
	ili.Cmd(ILI9341_DISPON)

	ili.Cmd(ILI9341_PIXFMT)
	ili.WriteByte(0x55) // 16 bit 565 format.

	ili.Cmd(ILI9341_MADCTL)
	ili.WriteByte(0x48) // Screen orientation.

	ili.Cmd(ILI9341_RAMWR)

	ili.SPI.P.SetWordSize(16)

	line := [240]uint16{100: 0xffff, 140: 0xffff}
	t := rtos.Nanosec()
	for i := 0; i < 320; i++ {
		ili.SPI.WriteRead16(line[:], nil)
	}
	fmt.Printf("%d ms\n", (rtos.Nanosec()-t+0.5e6)/1e6)
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
	irq.SPI3:          iliSPIISR,
	irq.DMA2_Channel1: iliRxDMAISR,
	irq.DMA2_Channel2: iliTxDMAISR,
}
