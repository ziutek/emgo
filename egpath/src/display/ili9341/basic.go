package ili9341

import (
	"image"
)

// DrawPoint draws a point (one pixel). 16-bit command.
func (a *Area) DrawPoint(p image.Point) {
	if !p.In(a.Bounds()) {
		return
	}
	p = p.Add(a.P0())
	dci := a.disp.dci // Reduces code size.
	dci.Cmd2(CASET)
	dci.WriteWord(uint16(p.X))
	dci.WriteWord(uint16(p.X))
	dci.Cmd2(PASET)
	dci.WriteWord(uint16(p.Y))
	dci.WriteWord(uint16(p.Y))
	dci.Cmd2(RAMWR)
	dci.WriteWord(uint16(a.color))
}

// rawFillRect helps to reduce code size (dci is an interface, that causes
// indirect method calls).
func (a *Area) rawFillRect(x0, y0, x1, y1, wxh int) {
	x0 += int(a.x0)
	y0 += int(a.y0)
	x1 += int(a.x0)
	y1 += int(a.y0)
	dci := a.disp.dci // Reduces code size.
	dci.Cmd2(CASET)
	dci.WriteWord(uint16(x0))
	dci.WriteWord(uint16(x1))
	dci.Cmd2(PASET)
	dci.WriteWord(uint16(y0))
	dci.WriteWord(uint16(y1))
	dci.Cmd2(RAMWR)
	dci.Fill(uint16(a.color), wxh)
}

// FillRect draws a filled rectangle. 16-bit command.
func (a *Area) FillRect(r image.Rectangle) {
	r = r.Canon().Intersect(a.Bounds())
	if !r.Empty() {
		a.rawFillRect(r.Min.X, r.Min.Y, r.Max.X-1, r.Max.Y-1, r.Dx()*r.Dy())
	}
}

func (a *Area) hline(x0, y0, x1 int) {
	r := a.Bounds()
	if y0 < r.Min.Y || y0 >= r.Max.Y {
		return
	}
	if x0 < r.Min.X {
		x0 = r.Min.X
	}
	if x1 >= r.Max.X {
		x1 = r.Max.X - 1
	}
	if x0 <= x1 {
		a.rawFillRect(x0, y0, x1, y0, x1-x0+1)
	}
}

func (a *Area) vline(x0, y0, y1 int) {
	r := a.Bounds()
	if x0 < r.Min.X || x0 >= r.Max.X {
		return
	}
	if y0 < r.Min.Y {
		y0 = r.Min.Y
	}
	if y1 >= r.Max.Y {
		y1 = r.Max.Y - 1
	}
	if y0 <= y1 {
		a.rawFillRect(x0, y0, x0, y1, y1-y0+1)
	}
}
