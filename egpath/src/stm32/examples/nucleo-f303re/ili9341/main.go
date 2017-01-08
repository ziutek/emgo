package main

import (
	"delay"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"rtos"

	"display/ili9341"

	"stm32/ilidci"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var dci *ilidci.DCI

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
	delay.Millisec(1) // Reset pulse.
	ilireset.Set()
	delay.Millisec(5) // Wait for reset.
	ilics.Clear()

	dci = ilidci.NewDCI(ilispi, ilidc)
}

func main() {
	delay.Millisec(100)
	spibus := dci.SPI().P.Bus()
	baudrate := dci.SPI().P.Baudrate(dci.SPI().P.Conf())
	fmt.Printf(
		"\nSPI on %s (%d MHz). SPI speed: %d bps.\n\n",
		spibus, spibus.Clock()/1e6, baudrate,
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

	ili := ili9341.NewDisplay(dci)

	width := ili.Bounds().Dx()
	height := ili.Bounds().Dy()
	wxh := width * height

	dci.Cmd(ILI9341_SLPOUT)
	delay.Millisec(120)
	dci.Cmd(ILI9341_DISPON)

	dci.Cmd(ILI9341_PIXFMT)
	dci.Byte(0x55) // 16 bit 565 format.

	dci.Cmd(ILI9341_MADCTL)
	dci.Byte(0xe8) // Screen orientation.

	dci.SetWordSize(16)

	dci.Cmd16(ILI9341_CASET)
	dci.Word(0)
	dci.Word(uint16(width - 1))
	dci.Cmd16(ILI9341_PASET)
	dci.Word(0)
	dci.Word(uint16(height - 1))

	dci.Cmd16(ILI9341_RAMWR)

	const N = 8
	start := rtos.Nanosec()
	for i := 0; i < N; i++ {
		dci.Fill(0xffff, wxh)
		dci.Fill(0, wxh)
	}
	fps := N * 2 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf(
		"dci.Fill speed: %.1f FPS (%.0f bps).\n", fps, fps*float32(wxh*16),
	)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		ili.SetColor(0xffff)
		ili.Rect(ili.Bounds())
		ili.SetColor(0)
		ili.Rect(ili.Bounds())
	}
	fps = N * 2 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf(
		"ili.Rect speed: %.1f FPS (%.0f bps).\n", fps, fps*float32(wxh*16),
	)

	var rnd rand.XorShift64
	rnd.Seed(rtos.Nanosec())

	p0 := image.Pt(-20, 10)
	p1 := image.Pt(350, 100)
	p2 := image.Pt(-20, 150)

	ili.SetColor(0xff00)
	ili.Point(p0)
	ili.Point(p1)
	ili.Point(p2)

	delay.Millisec(1000)

	ili.SetColor(0x00ff)
	ili.Line(p0, p1)
	ili.Line(p1, p2)
	ili.Line(p2, p0)

	p0 = image.Pt(100, -20)
	p1 = image.Pt(150, 250)
	p2 = image.Pt(180, -20)

	ili.SetColor(0xff00)
	ili.Point(p0)
	ili.Point(p1)
	ili.Point(p2)

	delay.Millisec(1000)

	ili.SetColor(0x00ff)
	ili.Line(p0, p1)
	ili.Line(p1, p2)
	ili.Line(p2, p0)

	ili.SetColor(0xf00f)
	ili.Line(image.Pt(-10, 120), image.Pt(350, 120))
	ili.Line(image.Pt(160, -10), image.Pt(160, 250))

	delay.Millisec(2000)

	for {
		v := rnd.Uint64()
		vh := uint32(v >> 32)
		vl := uint32(v)

		ili.SetColor(color.RGB16(vl))

		var r image.Rectangle

		vl >>= 16
		r.Min.Y = int(vl&0xff) - (256-height)/2
		vl >>= 8
		r.Max.Y = int(vl&0xff) - (256-height)/2

		r.Min.X = int(vh&0x1ff) - (512-width)/2
		vh >>= 9
		r.Max.X = int(vh&0x1ff) - (512-width)/2

		ili.Rect(r)
		//ili.Point(r.Min)
		//ili.Point(r.Max)
	}
}

func move(p image.Point) image.Point {
	switch {
	case p.Y == 40:
		if p.X += 5; p.X == 260 {
			p.Y += 5
		}
	case p.Y == 200:
		if p.X -= 5; p.X == 60 {
			p.Y -= 5
		}
	case p.X == 60:
		if p.Y -= 5; p.Y == 40 {
			p.X += 5
		}
	case p.X == 260:
		if p.Y += 5; p.Y == 200 {
			p.X -= 5
		}
	}
	return p
}

func iliSPIISR() {
	dci.SPI().ISR()
}

func iliRxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().RxDMA)
}

func iliTxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().TxDMA)
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
