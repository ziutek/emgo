package ili9341

import (
	"image"
	"image/color"
)

type Font struct {
}

// TextWriter allows to write a text on the display.
type TextWriter struct {
	area  *Area
	font  *Font
	color uint16
	pos   image.Point
}

func (a *Area) TextWriter(f *Font) TextWriter {
	return TextWriter{area: a, font: f}
}

func (w *TextWriter) SetPos(p image.Point) {
	w.pos = p
}

func (w *TextWriter) Pos() image.Point {
	return w.pos
}

func (w *TextWriter) SetColorRGB(r, g, b byte) {
	w.color = uint16(r)>>3<<11 | uint16(g)>>2<<5 | uint16(b)>>3
}

func (w *TextWriter) SetColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	w.color = uint16(r>>11<<11 | g>>10<<5 | b>>11)
}

func (w *TextWriter) WriteString(s string) (int, error) {
	return 0, nil
}

func (w *TextWriter) Write(s []byte) (int, error) {
	return 0, nil
}
