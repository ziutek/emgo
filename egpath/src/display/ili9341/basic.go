package ili9341

import (
	"image"
)

// Point draws a point (one pixel). 16-bit command.
func (d *Display) Point(p image.Point) {
	if !p.In(d.Bounds()) {
		return
	}
	dci := d.dci // Reduces code size.
	dci.Cmd16(CASET)
	dci.Word(uint16(p.X))
	dci.Word(uint16(p.X))
	dci.Cmd16(PASET)
	dci.Word(uint16(p.Y))
	dci.Word(uint16(p.Y))
	dci.Cmd16(RAMWR)
	dci.Word(uint16(d.color))
}

// rawRect is to reduce code size (dci is interface: indirect method calls).
func (d *Display) rawRect(x0, y0, x1, y1, wxh int) {
	dci := d.dci // Reduces code size.
	dci.Cmd16(CASET)
	dci.Word(uint16(x0))
	dci.Word(uint16(x1))
	dci.Cmd16(PASET)
	dci.Word(uint16(y0))
	dci.Word(uint16(y1))
	dci.Cmd16(RAMWR)
	dci.Fill(uint16(d.color), wxh)
}

// Rect draws a rectangle. 16-bit command.
func (d *Display) Rect(r image.Rectangle) {
	r = r.Canon().Intersect(d.Bounds())
	if !r.Empty() {
		d.rawRect(r.Min.X, r.Min.Y, r.Max.X-1, r.Max.Y-1, r.Dx()*r.Dy())
	}
}

func (d *Display) hline(x0, y0, x1 int) {
	r := d.Bounds()
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
		d.rawRect(x0, y0, x1, y0, x1-x0+1)
	}
}

func (d *Display) vline(x0, y0, y1 int) {
	r := d.Bounds()
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
		d.rawRect(x0, y0, x0, y1, y1-y0+1)
	}
}

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

// Line draws a line from p0 to p1 (including both points). 16-bit command.
func (d *Display) Line(p0, p1 image.Point) {
	dp := p1.Sub(p0)
	if dp.Y == 0 {
		if dp.X < 0 {
			p1.X, p0.X = p0.X, p1.X
		}
		d.hline(p0.X, p0.Y, p1.X)
		return
	}
	if dp.X == 0 {
		if dp.Y < 0 {
			p1.Y, p0.Y = p0.Y, p1.Y
		}
		d.vline(p0.X, p0.Y, p1.Y)
		return
	}
	vl := abs(dp.Y) > abs(dp.X)
	if vl {
		p0.X, p0.Y = p0.Y, p0.X
		p1.X, p1.Y = p1.Y, p1.X
	}
	if p0.X > p1.X {
		p0, p1 = p1, p0
	}
	dp = p1.Sub(p0).Mul(2)
	sy := 1
	if dp.Y < 0 {
		dp.Y = -dp.Y
		sy = -sy
	}
	e := p0.X - p1.X
	for x := p0.X; x <= p1.X; x++ {
		e += dp.Y
		if e > 0 {
			if vl {
				d.vline(p0.Y, p0.X, x)
			} else {
				d.hline(p0.X, p0.Y, x)
			}
			p0.X = x + 1
			p0.Y += sy
			e -= dp.X
		}
	}
	if p0.X <= p1.X {
		if vl {
			d.vline(p0.Y, p0.X, p1.X)
		} else {
			d.hline(p0.X, p0.Y, p1.X)
		}
	}
}

// Line1 draws a line from p0 to p1 (including both points). 16-bit command.
// Line1 uses less memory for code than Line but is generally slower (can be
// faster for very short lines: 1-3 pixels). Use Line1 if you are very short of
// Flash space and do not care about speed or to draw very short lines.
func (d *Display) Line1(p0, p1 image.Point) {
	dp := p1.Sub(p0)
	sx, sy := 1, 1
	if dp.X < 0 {
		sx = -sx
	}
	if dp.Y < 0 {
		sy = -sy
	}
	dp.X = abs(dp.X)
	dp.Y = abs(dp.Y)
	e := dp.X - dp.Y
	for {
		d.Point(p0)
		if p0 == p1 {
			return
		}
		e2 := 2 * e
		if e2 > -dp.Y {
			e -= dp.Y
			p0.X += sx
		}
		if e2 < dp.X {
			e += dp.X
			p0.Y += sy
		}
	}
}
