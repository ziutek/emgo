package ili9341

import (
	"image"
	"image/color"
)

type Area struct {
	disp   *Display
	rect   image.Rectangle
	x0, y0 uint16
	width  uint16
	height uint16
	color  uint16
	swapWH bool
}

func (a *Area) P0() image.Point {
	return image.Pt(int(a.x0), int(a.y0))
}

func (a *Area) updateBounds() {
	wh := a.rect.Intersect(a.disp.Bounds())
	a.x0 = uint16(wh.Min.X)
	a.y0 = uint16(wh.Min.Y)
	a.width = uint16(wh.Dx())
	a.height = uint16(wh.Dy())
	a.swapWH = a.disp.swapWH
}

func (a *Area) Bounds() image.Rectangle {
	if a.swapWH != a.disp.swapWH {
		a.updateBounds()
	}
	return image.Rectangle{Max: image.Pt(int(a.width), int(a.height))}
}

// SetColor sets the color used by drawing methods.
func (a *Area) SetColorRGB(r, g, b byte) {
	a.color = uint16(r)>>3<<11 | uint16(g)>>2<<5 | uint16(b)>>3
}

// SetColor sets the color used by drawing methods.
func (a *Area) SetColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	a.color = uint16(r>>11<<11 | g>>10<<5 | b>>11)
}
