package main

import (
	"delay"
	"fmt"
	"math/rand"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/ilidci"
)

var ili *ilidci.DCI

func init() {
	system.Setup(8, 1, 72/8)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7
	ilics := gpio.A.Pin(15)

	gpio.B.EnableClock(true)
	ilidc := gpio.B.Pin(7)

	gpio.C.EnableClock(true)
	//spiport, sck, miso, mosi := gpio.C, gpio.Pin10, gpio.Pin11, gpio.Pin12
	ilireset := gpio.C.Pin(13) // Max output: 2 MHz, 3 mA.

	// SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	//spiport.SetAltFunc(sck|miso|mosi, gpio.SPI3_AF6) // PC10-12
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1) // PA5-7
	//d := dma.DMA2
	d := dma.DMA1
	d.EnableClock(true)
	//ilispi := spi.NewDriver(spi.SPI3, d.Channel(1, 0), d.Channel(2, 0))
	ilispi := spi.NewDriver(spi.SPI1, d.Channel(2, 0), d.Channel(3, 0))
	ilispi.P.EnableClock(true)
	ilispi.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			ilispi.P.BR(36e6) | // 36 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	ilispi.P.SetWordSize(8)
	ilispi.P.Enable()
	//rtos.IRQ(irq.SPI3).Enable()
	//rtos.IRQ(irq.DMA2_Channel1).Enable()
	//rtos.IRQ(irq.DMA2_Channel2).Enable()
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
	ilireset.Clear()
	delay.Millisec(1) // Reset pulse.
	ilireset.Set()
	delay.Millisec(5) // Wait for reset.
	ilics.Clear()

	ili = ilidci.NewDCI(ilispi, ilidc)
}

func main() {
	delay.Millisec(100)
	spibus := ili.SPI().P.Bus()
	fmt.Printf(
		"\nSPI on %s (%d MHz). SPI speed: %d Hz.\n\n",
		spibus, spibus.Clock()/1e6, ili.SPI().P.Baudrate(ili.SPI().P.Conf()),
	)

	const (
		ILI9341_SLPOUT  = 0x11
		ILI9341_DISPOFF = 0x28
		ILI9341_DISPON  = 0x29
		ILI9341_RAMWR   = 0x2C
		ILI9341_MADCTL  = 0x36
		ILI9341_PIXFMT  = 0x3A
		ILI9341_CASET   = 0x2A
		ILI9341_PASET   = 0x2B
	)

	ili.Cmd(ILI9341_SLPOUT)
	delay.Millisec(120)
	ili.Cmd(ILI9341_DISPON)

	ili.Cmd(ILI9341_PIXFMT)
	ili.Byte(0x55) // 16 bit 565 format.

	ili.Cmd(ILI9341_MADCTL)
	ili.Byte(0xe8) // Screen orientation.

	ili.SetWordSize(16)

	ili.Cmd16(ILI9341_CASET)
	ili.Word(0)
	ili.Word(320 - 1)
	ili.Cmd16(ILI9341_PASET)
	ili.Word(0)
	ili.Word(240 - 1)

	ili.Cmd16(ILI9341_RAMWR)
	ili.Fill(0, 320*240)

	var rnd rand.XorShift64
	rnd.Seed(rtos.Nanosec())

	for {
		v := rnd.Uint64()
		vh := uint32(v >> 32)
		vl := uint32(v)
		c := uint16(vl)
		vl >>= 16
		x := vl & 0xff
		vl >>= 8
		y := vl & 0x7f
		w := vh&0x7f + 32
		vh >>= 8
		h := vh&0x3f + 64
		if x+w > 320 {
			w = 320 - x
		}
		if y+h > 240 {
			h = 240 - y
		}

		ili.Cmd16(ILI9341_CASET)
		ili.Word(uint16(x))
		ili.Word(uint16(x + w - 1))

		ili.Cmd16(ILI9341_PASET)
		ili.Word(uint16(y))
		ili.Word(uint16(y + h - 1))

		ili.Cmd16(ILI9341_RAMWR)
		ili.Fill(c, int(w*h))
	}
}

func iliSPIISR() {
	ili.SPI().ISR()
}

func iliRxDMAISR() {
	ili.SPI().DMAISR(ili.SPI().RxDMA)
}

func iliTxDMAISR() {
	ili.SPI().DMAISR(ili.SPI().TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	//irq.SPI3: iliSPIISR,
	//irq.DMA2_Channel1: iliRxDMAISR,
	//irq.DMA2_Channel2: iliTxDMAISR,
	irq.SPI1:          iliSPIISR,
	irq.DMA1_Channel2: iliRxDMAISR,
	irq.DMA1_Channel3: iliTxDMAISR,
}
