package ili9341

import (
	"image"
	"image/color"
)

type Display struct {
	dci   DCI
	color color.RGB16
	w, h  uint16
}

// MakeDisplay returns initialised Display value.
func MakeDisplay(dci DCI) Display {
	return Display{dci: dci, w: 240, h: 320}
}

// NewDisplay works like MakeDisplay but returns a pointer to the heap allocated
// variable.
func NewDisplay(dci DCI) *Display {
	d := new(Display)
	*d = MakeDisplay(dci)
	return d
}

// DCI allows to direct access to the internal DCI.
func (d *Display) DCI() DCI {
	return d.dci
}

// Err returns and clears internal error variable.
func (d *Display) Err() error {
	return d.dci.Err()
}

// Bounds current display bounds.
func (d *Display) Bounds() image.Rectangle {
	return image.Rectangle{Max: image.Pt(int(d.w), int(d.h))}
}

// SetWordSize changes the data word size.
func (d *Display) SetWordSize(size int) {
	d.dci.SetWordSize(size)
}

// SetColor sets a color used by drawing methods.
func (d *Display) SetColor(c color.RGB16) {
	d.color = c
}
