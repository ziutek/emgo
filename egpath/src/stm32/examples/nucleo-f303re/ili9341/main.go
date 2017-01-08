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

var (
	lcdspi *spi.Driver
	lcd    *ili9341.Display
)

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

	lcd = ili9341.NewDisplay(ilidci.NewDCI(lcdspi, ilidc))
}

func main() {
	delay.Millisec(100)
	spibus := lcdspi.P.Bus()
	baudrate := lcdspi.P.Baudrate(lcdspi.P.Conf())
	fmt.Printf(
		"\nSPI on %s (%d MHz). SPI speed: %d bps.\n\n",
		spibus, spibus.Clock()/1e6, baudrate,
	)

	width := lcd.Bounds().Dx()
	height := lcd.Bounds().Dy()
	wxh := width * height

	lcd.SlpOut()
	delay.Millisec(120)
	lcd.DispOn()
	lcd.PixSet(ili9341.PF16) // 16-bit pixel format.
	lcd.MADCtl(ili9341.MY | ili9341.MX | ili9341.MV | ili9341.BGR)
	lcd.SetWordSize(16)

	dci := lcd.DCI()
	dci.Cmd16(ili9341.CASET)
	dci.Word(0)
	dci.Word(uint16(width - 1))
	dci.Cmd16(ili9341.PASET)
	dci.Word(0)
	dci.Word(uint16(height - 1))
	dci.Cmd16(ili9341.RAMWR)

	const N = 10
	start := rtos.Nanosec()
	for i := 0; i < N; i++ {
		dci.Fill(0xffff, wxh)
		dci.Fill(0, wxh)
	}
	fps := N * 2 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf(
		"dci.Fill  speed: %4.1f fps (%.0f bps).\n", fps, fps*float32(wxh*16),
	)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		lcd.SetColor(0xffff)
		lcd.Rect(lcd.Bounds())
		lcd.SetColor(0)
		lcd.Rect(lcd.Bounds())
	}
	fps = N * 2 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf(
		"lcd.Rect  speed: %4.1f fps (%.0f bps).\n", fps, fps*float32(wxh*16),
	)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		lcd.SetColor(0xffff)
		lcd.Line(image.Pt(0, 0), image.Pt(319, 239))
		lcd.Line(image.Pt(0, 239), image.Pt(319, 0))
		lcd.Line(image.Pt(-10, 120), image.Pt(350, 120))
		lcd.Line(image.Pt(160, -10), image.Pt(160, 250))
		lcd.SetColor(0)
		lcd.Line(image.Pt(0, 0), image.Pt(319, 239))
		lcd.Line(image.Pt(0, 239), image.Pt(319, 0))
		lcd.Line(image.Pt(-10, 120), image.Pt(350, 120))
		lcd.Line(image.Pt(160, -10), image.Pt(160, 250))
	}
	lps := N * 8 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("lcd.Line  speed: %4.0f lps.\n", lps)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		lcd.SetColor(0xffff)
		lcd.Line1(image.Pt(0, 0), image.Pt(319, 239))
		lcd.Line1(image.Pt(0, 239), image.Pt(319, 0))
		lcd.Line1(image.Pt(-10, 120), image.Pt(350, 120))
		lcd.Line1(image.Pt(160, -10), image.Pt(160, 250))
		lcd.SetColor(0)
		lcd.Line1(image.Pt(0, 0), image.Pt(319, 239))
		lcd.Line1(image.Pt(0, 239), image.Pt(319, 0))
		lcd.Line1(image.Pt(-10, 120), image.Pt(350, 120))
		lcd.Line1(image.Pt(160, -10), image.Pt(160, 250))
	}
	lps = N * 8 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("lcd.Line1 speed: %4.0f lps.\n", lps)

	p0 := image.Pt(40, 40)
	p1 := image.Pt(200, 100)
	p2 := image.Pt(60, 150)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		lcd.SetColor(0xffff)
		lcd.Line(p0, p1)
		lcd.Line(p1, p2)
		lcd.Line(p2, p0)
		lcd.SetColor(0)
		lcd.Line(p0, p1)
		lcd.Line(p1, p2)
		lcd.Line(p2, p0)
	}
	lps = N * 6 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("lcd.Line  speed: %4.0f lps.\n", lps)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		lcd.SetColor(0xffff)
		lcd.Line1(p0, p1)
		lcd.Line1(p1, p2)
		lcd.Line1(p2, p0)
		lcd.SetColor(0)
		lcd.Line1(p0, p1)
		lcd.Line1(p1, p2)
		lcd.Line1(p2, p0)
	}
	lps = N * 6 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("lcd.Line1 speed: %4.0f lps.\n", lps)

	delay.Millisec(1e3)

	var rnd rand.XorShift64
	rnd.Seed(rtos.Nanosec())

	for {
		v := rnd.Uint64()
		vh := uint32(v >> 32)
		vl := uint32(v)

		lcd.SetColor(color.RGB16(vl))

		var r image.Rectangle

		vl >>= 16
		r.Min.Y = int(vl&0xff) - (256-height)/2
		vl >>= 8
		r.Max.Y = int(vl&0xff) - (256-height)/2

		r.Min.X = int(vh&0x1ff) - (512-width)/2
		vh >>= 9
		r.Max.X = int(vh&0x1ff) - (512-width)/2

		lcd.Rect(r)
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
