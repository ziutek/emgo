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
	color color.RGB16
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

func (w *TextWriter) SetColor(c color.RGB16) {
	w.color = c
}

func (w *TextWriter) WriteString(s string) (int, error) {
	return 0, nil
}

func (w *TextWriter) Write(s []byte) (int, error) {
	return 0, nil
}
