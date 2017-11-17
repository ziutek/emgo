package ili9341test

import (
	"delay"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"rtos"

	"display/ili9341"
)

func perSec(n int, start int64) float32 {
	return float32(n) * 1e9 / float32(rtos.Nanosec()-start)
}

func printSpeed(what string, ps float32) {
	fmt.Printf("%s speed: %4.1f/s.\n", what, ps)
}

func Run(lcd *ili9341.Display, n int, init bool) {
	fmt.Printf("TEST BEGIN\n")

	if init {
		lcd.SlpOut()
		delay.Millisec(120)
		lcd.DispOn()
		lcd.PixSet(ili9341.PF16) // 16-bit pixel format.
		lcd.MADCtl(ili9341.MY | ili9341.MX | ili9341.MV | ili9341.BGR)
	}
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

	start := rtos.Nanosec()
	for i := 0; i < n; i++ {
		dci.Fill(0x0fff, wxh)
		dci.Fill(0, wxh)
	}
	fps := perSec(n*2, start)
	fmt.Printf(
		"dci.Fill          speed: %4.1f fps (%.0f bps).\n",
		fps, fps*float32(wxh*16),
	)

	scr := lcd.Area(lcd.Bounds())

	start = rtos.Nanosec()
	for i := n; i > 0; i-- {
		scr.SetColor(0xfff0)
		scr.FillRect(scr.Bounds())
		scr.SetColor(0)
		scr.FillRect(scr.Bounds())
	}
	fps = perSec(n*2, start)
	fmt.Printf(
		"scr.FillRect      speed: %4.1f fps (%.0f bps).\n",
		fps, fps*float32(wxh*16),
	)

	start = rtos.Nanosec()
	for i := n; i > 0; i-- {
		scr.SetColor(0xf0ff)
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
	printSpeed("scr.DrawLine     ", perSec(n*8, start))

	start = rtos.Nanosec()
	for i := 0; i < n; i++ {
		scr.SetColor(0xff0f)
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
	printSpeed("scr.DrawLine_    ", perSec(n*8, start))

	p0 := image.Pt(40, 40)
	p1 := image.Pt(300, 100)
	p2 := image.Pt(60, 200)

	start = rtos.Nanosec()
	for i := 0; i < n; i++ {
		scr.SetColor(0xf0ff)
		scr.DrawLine(p0, p1)
		scr.DrawLine(p1, p2)
		scr.DrawLine(p2, p0)
		scr.SetColor(0)
		scr.DrawLine(p0, p1)
		scr.DrawLine(p1, p2)
		scr.DrawLine(p2, p0)
	}
	printSpeed("scr.DrawLine     ", perSec(n*6, start))

	start = rtos.Nanosec()
	for i := 0; i < n; i++ {
		scr.SetColor(0xff0f)
		scr.DrawLine_(p0, p1)
		scr.DrawLine_(p1, p2)
		scr.DrawLine_(p2, p0)
		scr.SetColor(0)
		scr.DrawLine_(p0, p1)
		scr.DrawLine_(p1, p2)
		scr.DrawLine_(p2, p0)
	}
	printSpeed("scr.DrawLine_    ", perSec(n*6, start))

	p0 = scr.Bounds().Max.Div(2)
	r := p0.X
	if r > p0.Y {
		r = p0.Y
	}
	r--
	start = rtos.Nanosec()
	for i := 0; i < n; i++ {
		scr.SetColor(0xffff)
		scr.FillCircle(p0, r)
		scr.SetColor(0)
		scr.FillCircle(p0, r)
	}
	printSpeed("scr.FillCircle   ", perSec(n*2, start))

	start = rtos.Nanosec()
	for i := 0; i < n; i++ {
		scr.SetColor(0xffff)
		scr.DrawCircle(p0, r)
		scr.SetColor(0)
		scr.DrawCircle(p0, r)
	}
	printSpeed("scr.DrawCircle   ", perSec(n*2, start))

	var rnd rand.XorShift64
	rnd.Seed(1)
	var dt int64
	for i := n * 100; i > 0; i-- {
		v := rnd.Uint64()
		vl := uint32(v)
		vh := uint32(v >> 32)

		r := scr.Bounds()
		r.Min.Y += int(vh & 0xff)
		r.Max.Y -= int(vh >> 8 & 0xff)
		r.Min.X += int(vh >> 16 & 0xff)
		r.Max.X -= int(vh >> 24 & 0xff)

		start := rtos.Nanosec()
		scr.SetColor(color.RGB16(vl))
		if vl>>16&3 != 0 {
			scr.FillRect(r)
		} else {
			scr.FillCircle(r.Min, r.Max.X/4)
		}
		dt += rtos.Nanosec() - start
	}
	printSpeed("Rand. rect/circle", float32(n*100)*1e9/float32(dt))

	fmt.Printf("TEST END\n")
}
