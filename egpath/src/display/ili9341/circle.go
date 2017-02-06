package ili9341

import (
	"image"
)

// DrawCircle draws empty circle. 16-bit command.
func (a *Area) DrawCircle(p0 image.Point, r int) {
	x, y, e := r, 0, 1-r
	for x >= y {
		a.DrawPoint(p0.Add(image.Pt(-x, y)))
		a.DrawPoint(p0.Add(image.Pt(x, y)))
		a.DrawPoint(p0.Add(image.Pt(-x, -y)))
		a.DrawPoint(p0.Add(image.Pt(x, -y)))
		a.DrawPoint(p0.Add(image.Pt(-y, x)))
		a.DrawPoint(p0.Add(image.Pt(y, x)))
		a.DrawPoint(p0.Add(image.Pt(-y, -x)))
		a.DrawPoint(p0.Add(image.Pt(y, -x)))
		y++
		e += 2*y + 1
		if e > 0 {
			x--
			e -= 2 * x
		}
	}
}

// FillCircle draws filled circle. 16-bit command.
func (a *Area) FillCircle(p0 image.Point, r int) {
	// Fill four sides.
	x, y, e := r, 0, 1-r
	for x > y {
		e += 2*y + 3
		if e > 0 {
			y0, y1 := p0.Y-y, p0.Y+y
			a.vline(p0.X+x, y0, y1)
			a.vline(p0.X-x, y0, y1)
			x0, x1 := p0.X-y, p0.X+y
			a.hline(x0, p0.Y-x, x1)
			a.hline(x0, p0.Y+x, x1)
			x--
			e -= 2 * x
		}
		y++
	}
	// Fill central rectangle.
	a.FillRect(image.Rectangle{
		p0.Sub(image.Pt(x, y)), p0.Add(image.Pt(x+1, y+1)),
	})
}
