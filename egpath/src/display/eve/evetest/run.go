package evetest

import (
	"delay"
	"errors"
	"fmt"
	"time"

	"display/eve"
)

func waitTouch(lcd *eve.Driver) {
	delay.Millisec(100)
	lcd.ClearIntFlags(eve.INT_TOUCH)
	lcd.Wait(eve.INT_TOUCH)
}

func Run(lcd *eve.Driver) error {
	lcd.SetBacklight(64)

	width, height := lcd.Width(), lcd.Height()

	dl := lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.Vertex2f(200<<4, 100<<4)
	dl.Display()

	lcd.SwapDL()
	waitTouch(lcd)

	dl = lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.PointSize(70 << 4)
	dl.Vertex2f(200<<4, 100<<4)
	dl.Display()

	lcd.SwapDL()
	waitTouch(lcd)

	dl = lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.PointSize(70 << 4)
	dl.Vertex2f(200<<4, 100<<4)
	dl.ColorRGB(0x0000FF)
	dl.PointSize(50 << 4)
	dl.Vertex2f(240<<4, 150<<4)
	dl.Display()

	lcd.SwapDL()
	waitTouch(lcd)

	dl = lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.PointSize(70 << 4)
	dl.Vertex2f(200<<4, 100<<4)
	dl.ColorRGB(0x0000FF)
	dl.ColorA(128)
	dl.PointSize(50 << 4)
	dl.Vertex2f(240<<4, 150<<4)
	dl.Display()

	lcd.SwapDL()
	waitTouch(lcd)

	dl = lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.BITMAPS)
	dl.BitmapHandle(31)
	dl.Cell('E')
	dl.Vertex2f(200<<4, 100<<4)
	dl.Display()

	lcd.SwapDL()
	waitTouch(lcd)

	dl = lcd.DL(-1)
	dl.Clear(eve.CST)
	dl.Begin(eve.BITMAPS)
	dl.BitmapHandle(31)
	dl.Cell('E')
	dl.Vertex2f(200<<4, 100<<4)
	dl.Cell('V')
	dl.Vertex2f(224<<4, 100<<4)
	dl.Cell('E')
	dl.Vertex2f(250<<4, 100<<4)
	dl.Display()

	lcd.SwapDL()
	waitTouch(lcd)

	ge := lcd.GE(-1)
	ge.DLStart()
	ge.Clear(eve.CST)
	ge.TextString(width/2, height/2, 31, eve.OPT_CENTER, "Hello world!")
	ge.Display()
	ge.Swap()
	lcd.Wait(eve.INT_CMDEMPTY)

	waitTouch(lcd)

	lcd.Write(0, gopherMask[:])
	addr := (len(gopherMask) + 3) &^ 3

	ge.DLStart()
	ge.LoadImageBytes(addr, eve.OPT_RGB565, gopher[:])
	lcd.Wait(eve.INT_CMDEMPTY) // A lot of data sent. Ensure free space.
	ge.BitmapHandle(1)
	ge.BitmapLayout(eve.L1, 216/8, 251)
	ge.BitmapSize(eve.DEFAULT, 211, 251)
	ge.Clear(eve.CST)
	ge.Gradient(0, 0, 0x001155, 0, height, 0x772200)
	ge.Begin(eve.BITMAPS)
	ge.ColorMask(eve.A)
	ge.Clear(eve.C)
	ge.BitmapHandle(1)
	ge.Vertex2f(31*16, 21*16)
	ge.ColorMask(eve.RGBA)
	ge.BlendFunc(eve.DST_ALPHA, eve.ONE_MINUS_DST_ALPHA)
	ge.BitmapHandle(0)
	ge.Vertex2f(0, 0)
	ge.Display()
	ge.Swap()
	lcd.Wait(eve.INT_CMDEMPTY)

	waitTouch(lcd)
	delay.Millisec(200)

	ge.DLStart()
	ge.Clear(eve.CST)
	ge.TextString(
		width/2, height/2, 30, eve.OPT_CENTER,
		"Touch panel calibration",
	)
	addr = ge.Calibrate()
	lcd.Wait(eve.INT_CMDEMPTY)
	if lcd.ReadInt(addr) == 0 {
		return errors.New("touch calibration failed")
	}

	const button = 1
	for {
		tag := lcd.TouchTag()
		x, y := lcd.TouchScreenXY()
		ge.DLStart()
		ge.ClearColorRGB(0xc3a6f4)
		ge.Clear(eve.CST)
		ge.Gradient(0, 0, 0x0004ff, 0, height, 0xe08484)
		ge.Text(width-180, 20, 26, eve.DEFAULT)
		fmt.Fprintf(&ge, "x=%d y=%d tag=%d\000", x, y, tag)
		ge.Align32()
		ge.TextString(width/2, height/2, 30, eve.OPT_CENTER, "Hello World!")
		ge.Begin(eve.RECTS)
		ge.ColorA(128)
		ge.ColorRGB(0xFF8000)
		ge.Vertex2ii(260, 100, 0, 0)
		ge.Vertex2ii(360, 200, 0, 0)
		ge.ColorRGB(0x0080FF)
		ge.Vertex2ii(300, 160, 0, 0)
		ge.Vertex2ii(400, 260, 0, 0)
		ge.ColorRGB(0xFFFFFF)
		ge.ColorA(200)
		t := time.Now()
		h, m, s := t.Clock()
		ms := int(t.Nanosecond() / 1e6)
		ge.Clock(100, 100, 70, eve.OPT_NOBACK, h, m, s, ms)
		ge.ColorA(255)
		ge.Tag(button)
		buttonFont := byte(27)
		buttonStyle := uint16(eve.DEFAULT)
		if tag == button {
			buttonFont--
			buttonStyle |= eve.OPT_FLAT
			ge.TextString(300, height-70, 29, eve.DEFAULT, "Thanks!")
		}
		ge.ButtonString(
			40, height-70, 100, 40, buttonFont, buttonStyle,
			"Push me!",
		)
		ge.Display()
		ge.Swap()
		lcd.Wait(eve.INT_CMDEMPTY) // Wait for end of Swap (next frame).
	}

	return nil
}

/*
	var rnd rand.XorShift64
	rnd.Seed(1)

	const testLen = 120

	// Create display list directly in RAM_DL.

	n := 800 // Greather numbers can exceed the render limit (8192 pixel/line).
	fmt.Printf("\n%d bitmaps:", n)

	addr := 0
	lcd.W(addr).Write(LenaFaceRGB[:])
	addr += len(LenaFaceRGB)

	t := rtos.Nanosec()
	for i := 0; i < testLen; i++ {
		dl := lcd.DL(-1)
		dl.BitmapLayout(eve.RGB565, 80, 40)
		dl.BitmapSize(eve.DEFAULT, 40, 40)
		dl.Clear(eve.CST)
		dl.Begin(eve.BITMAPS)
		for k := 0; k < n; k++ {
			v := rnd.Uint32()
			c := v&0xFFFFFF | 0x808080
			x := int(v>>12) % lcd.Width()
			y := int(v>>23) % lcd.Height()
			dl.ColorRGB(c)
			dl.Vertex2f((x-20)*16, (y-20)*16)
		}
		dl.Display()
		lcd.SwapDL()
	}
	fmt.Printf(" %.2f fps.\n", testLen*1e9/float32(rtos.Nanosec()-t))
*/
