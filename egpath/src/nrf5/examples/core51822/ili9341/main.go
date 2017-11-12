package main

import (
	"delay"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"rtos"

	"display/ili9341"

	"nrf5/ilidci"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/spi"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/uart"
)

var (
	leds   [5]gpio.Pin
	u      *uart.Driver
	lcdspi *spi.Driver
	lcd    *ili9341.Display
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// GPIO

	p0 := gpio.P0

	ilireset := p0.Pin(0)
	ilidc := p0.Pin(1)
	ilimosi := p0.Pin(2)
	ilisck := p0.Pin(3)
	ilimiso := p0.Pin(4)

	utx := p0.Pin(9)
	urx := p0.Pin(11)

	for i := range leds {
		leds[i] = p0.Pin(18 + i)
	}

	ilicsn := p0.Pin(30)

	// LEDs

	for _, led := range leds {
		led.Setup(gpio.ModeOut)
	}

	// UART

	u = uart.NewDriver(uart.UART0, make([]byte, 80))
	u.P.StorePSEL(uart.RXD, urx)
	u.P.StorePSEL(uart.TXD, utx)
	u.P.StoreBAUDRATE(uart.Baud115200)
	u.Enable()
	//u.EnableRx()
	u.EnableTx()
	rtos.IRQ(u.P.NVIC()).Enable()
	fmt.DefaultWriter = u

	// LCD SPI

	lcdspi = spi.NewDriver(spi.SPI0)
	lcdspi.P.StorePSEL(spi.SCK, ilisck)
	lcdspi.P.StorePSEL(spi.MISO, ilimiso)
	lcdspi.P.StorePSEL(spi.MOSI, ilimosi)
	lcdspi.P.StoreFREQUENCY(spi.Freq8M)
	lcdspi.Enable()
	rtos.IRQ(lcdspi.P.NVIC()).Enable()

	// LCD controll

	p0.Setup(ilicsn.Mask()|ilidc.Mask()|ilireset.Mask(), gpio.ModeOut)
	ilicsn.Set()
	delay.Millisec(1) // Reset pulse.
	ilireset.Set()
	delay.Millisec(5) // Wait for reset.
	ilicsn.Clear()

	lcd = ili9341.NewDisplay(ilidci.NewDCI(lcdspi, ilidc))
}

func main() {
	lcd.SlpOut()
	delay.Millisec(120)
	lcd.DispOn()
	lcd.PixSet(ili9341.PF16) // 16-bit pixel format.
	lcd.MADCtl(ili9341.MY | ili9341.MX | ili9341.MV | ili9341.BGR)

	lcd.SetWordSize(16)

	width := lcd.Bounds().Dx()
	height := lcd.Bounds().Dy()
	wxh := width * height

	dci := lcd.DCI()
	dci.Cmd2(ili9341.CASET)
	dci.WriteWord(0)
	dci.WriteWord(uint16(width - 1))
	dci.Cmd2(ili9341.PASET)
	dci.WriteWord(0)
	dci.WriteWord(uint16(height - 1))
	dci.Cmd2(ili9341.RAMWR)

	const N = 4
	start := rtos.Nanosec()
	for i := 0; i < N; i++ {
		dci.Fill(0xffff, wxh)
		dci.Fill(0, wxh)
	}
	fps := N * 2 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf(
		"\r\n\r\ndci.Fill       speed: %4.1f fps (%.0f bps).\r\n",
		fps, fps*float32(wxh*16),
	)

	scr := lcd.Area(lcd.Bounds())

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		scr.SetColor(0xffff)
		scr.FillRect(scr.Bounds())
		scr.SetColor(0)
		scr.FillRect(scr.Bounds())
	}
	fps = N * 2 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf(
		"scr.FillRect   speed: %4.1f fps (%.0f bps).\r\n",
		fps, fps*float32(wxh*16),
	)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		scr.SetColor(0xffff)
		scr.DrawLine(image.Pt(0, 0), image.Pt(319, 239))
		scr.DrawLine(image.Pt(0, 239), image.Pt(319, 0))
		scr.DrawLine(image.Pt(-10, 120), image.Pt(350, 120))
		scr.DrawLine(image.Pt(160, -10), image.Pt(160, 250))
		scr.SetColor(0)
		scr.DrawLine(image.Pt(0, 0), image.Pt(319, 239))
		scr.DrawLine(image.Pt(0, 239), image.Pt(319, 0))
		scr.DrawLine(image.Pt(-10, 120), image.Pt(350, 120))
		scr.DrawLine(image.Pt(160, -10), image.Pt(160, 250))
	}
	lps := N * 8 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("scr.DrawLine   speed: %4.0f lps.\r\n", lps)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		scr.SetColor(0xffff)
		scr.DrawLine_(image.Pt(0, 0), image.Pt(319, 239))
		scr.DrawLine_(image.Pt(0, 239), image.Pt(319, 0))
		scr.DrawLine_(image.Pt(-10, 120), image.Pt(350, 120))
		scr.DrawLine_(image.Pt(160, -10), image.Pt(160, 250))
		scr.SetColor(0)
		scr.DrawLine_(image.Pt(0, 0), image.Pt(319, 239))
		scr.DrawLine_(image.Pt(0, 239), image.Pt(319, 0))
		scr.DrawLine_(image.Pt(-10, 120), image.Pt(350, 120))
		scr.DrawLine_(image.Pt(160, -10), image.Pt(160, 250))
	}
	lps = N * 8 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("scr.DrawLine_  speed: %4.0f lps.\r\n", lps)

	p0 := image.Pt(40, 40)
	p1 := image.Pt(200, 100)
	p2 := image.Pt(60, 150)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		scr.SetColor(0xffff)
		scr.DrawLine(p0, p1)
		scr.DrawLine(p1, p2)
		scr.DrawLine(p2, p0)
		scr.SetColor(0)
		scr.DrawLine(p0, p1)
		scr.DrawLine(p1, p2)
		scr.DrawLine(p2, p0)
	}
	lps = N * 6 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("scr.DrawLine   speed: %4.0f lps.\r\n", lps)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		scr.SetColor(0xffff)
		scr.DrawLine_(p0, p1)
		scr.DrawLine_(p1, p2)
		scr.DrawLine_(p2, p0)
		scr.SetColor(0)
		scr.DrawLine_(p0, p1)
		scr.DrawLine_(p1, p2)
		scr.DrawLine_(p2, p0)
	}
	lps = N * 6 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("scr.DrawLine_  speed: %4.0f lps.\r\n", lps)

	p0 = scr.Bounds().Max.Div(2)
	r := p0.X
	if r > p0.Y {
		r = p0.Y
	}
	r--
	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		scr.SetColor(0xffff)
		scr.FillCircle(p0, r)
		scr.SetColor(0)
		scr.FillCircle(p0, r)
	}
	cps := N * 2 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("scr.FillCircle speed: %4.0f cps.\r\n", cps)

	start = rtos.Nanosec()
	for i := 0; i < N; i++ {
		scr.SetColor(0xffff)
		scr.DrawCircle(p0, r)
		scr.SetColor(0)
		scr.DrawCircle(p0, r)
	}
	cps = N * 2 * 1e9 / float32(rtos.Nanosec()-start)
	fmt.Printf("scr.DrawCircle speed: %4.0f cps.\r\n", cps)

	var rnd rand.XorShift64
	rnd.Seed(rtos.Nanosec())

	for {
		v := rnd.Uint64()
		vl := uint32(v)
		vh := uint32(v >> 32)

		r := scr.Bounds()
		r.Min.Y += int(vh & 0xff)
		r.Max.Y -= int(vh >> 8 & 0xff)
		r.Min.X += int(vh >> 16 & 0xff)
		r.Max.X -= int(vh >> 24 & 0xff)

		scr.SetColor(color.RGB16(vl))
		if vl>>16&3 != 0 {
			scr.FillRect(r)
		} else {
			scr.FillCircle(r.Min, r.Max.X/4)
		}
	}
}

func spiISR() {
	lcdspi.ISR()
}

func uartISR() {
	u.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:      rtcst.ISR,
	irq.SPI0_TWI0: spiISR,
	irq.UART0:     uartISR,
}
