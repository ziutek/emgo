package ili9341

import (
	"image"
)

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

// DrawLine draws a line from p0 to p1 (including both points). 16-bit command.
func (a *Area) DrawLine(p0, p1 image.Point) {
	dp := p1.Sub(p0)
	if dp.Y == 0 {
		if dp.X < 0 {
			p1.X, p0.X = p0.X, p1.X
		}
		a.hline(p0.X, p0.Y, p1.X)
		return
	}
	if dp.X == 0 {
		if dp.Y < 0 {
			p1.Y, p0.Y = p0.Y, p1.Y
		}
		a.vline(p0.X, p0.Y, p1.Y)
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
				a.vline(p0.Y, p0.X, x)
			} else {
				a.hline(p0.X, p0.Y, x)
			}
			p0.X = x + 1
			p0.Y += sy
			e -= dp.X
		}
	}
	if p0.X <= p1.X {
		if vl {
			a.vline(p0.Y, p0.X, p1.X)
		} else {
			a.hline(p0.X, p0.Y, p1.X)
		}
	}
}

// DrawLine_ draws a line from p0 to p1 (including both points). 16-bit
// command. DrawLine_ uses less memory for code than DrawLine but is generally
// slower (can be faster for very short lines: 1-3 pixels). Use DrawLine_ if
// you are very short of Flash space and do not care about speed or to draw
// very short lines.
func (a *Area) DrawLine_(p0, p1 image.Point) {
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
		a.DrawPoint(p0)
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
