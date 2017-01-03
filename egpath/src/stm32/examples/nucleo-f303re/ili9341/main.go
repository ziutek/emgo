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
)

var ili ILI9341

func init() {
	system.Setup(8, 1, 72/8)
	systick.Setup()

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7
	ili.cs = gpio.A.Pin(15)

	gpio.B.EnableClock(true)
	ili.dc = gpio.B.Pin(7)

	gpio.C.EnableClock(true)
	//spiport, sck, miso, mosi := gpio.C, gpio.Pin10, gpio.Pin11, gpio.Pin12
	ili.reset = gpio.C.Pin(13) // Max output: 2 MHz, 3 mA.

	// SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	//spiport.SetAltFunc(sck|miso|mosi, gpio.SPI3_AF6) // PC10-12
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1) // PA5-7
	//d := dma.DMA2
	d := dma.DMA1
	d.EnableClock(true)
	//ili.spi = spi.NewDriver(spi.SPI3, d.Channel(1, 0), d.Channel(2, 0))
	ili.spi = spi.NewDriver(spi.SPI1, d.Channel(2, 0), d.Channel(3, 0))
	ili.spi.P.EnableClock(true)
	ili.spi.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			ili.spi.P.BR(36e6) | // 36 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	ili.spi.P.SetWordSize(8)
	ili.spi.P.Enable()
	//rtos.IRQ(irq.SPI3).Enable()
	//rtos.IRQ(irq.DMA2_Channel1).Enable()
	//rtos.IRQ(irq.DMA2_Channel2).Enable()
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	// Controll

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.High}
	ili.cs.Setup(&cfg)
	ili.dc.Setup(&cfg)
	cfg.Speed = gpio.Low
	ili.reset.Setup(&cfg)
}

func main() {
	delay.Millisec(100)
	spibus := ili.spi.P.Bus()
	fmt.Printf(
		"\nSPI on %s (%d MHz). SPI speed: %d Hz.\n\n",
		spibus, spibus.Clock()/1e6, ili.spi.P.Baudrate(ili.spi.P.Conf()),
	)

	ili.Reset()
	delay.Millisec(5)
	ili.Select()
	ili.Cmd(ILI9341_SLPOUT)
	delay.Millisec(120)
	ili.Cmd(ILI9341_DISPON)

	ili.Cmd(ILI9341_PIXFMT)
	ili.WriteByte(0x55) // 16 bit 565 format.

	ili.Cmd(ILI9341_MADCTL)
	ili.WriteByte(0xe8) // Screen orientation.

	ili.spi.P.SetWordSize(16)

	ili.Cmd16(ILI9341_CASET)
	ili.WriteWord16(0)
	ili.WriteWord16(320 - 1)
	ili.Cmd16(ILI9341_PASET)
	ili.WriteWord16(0)
	ili.WriteWord16(240 - 1)

	ili.Cmd16(ILI9341_RAMWR)
	ili.Fill16(0, 320*240)

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
		ili.WriteWord16(uint16(x))
		ili.WriteWord16(uint16(x + w - 1))

		ili.Cmd16(ILI9341_PASET)
		ili.WriteWord16(uint16(y))
		ili.WriteWord16(uint16(y + h - 1))

		ili.Cmd16(ILI9341_RAMWR)
		ili.Fill16(c, int(w*h))
	}
}

func iliSPIISR() {
	ili.spi.ISR()
}

func iliRxDMAISR() {
	ili.spi.DMAISR(ili.spi.RxDMA)
}

func iliTxDMAISR() {
	ili.spi.DMAISR(ili.spi.TxDMA)
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
