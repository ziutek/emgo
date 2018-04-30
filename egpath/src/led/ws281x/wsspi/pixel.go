// Package wsspi allows to use SPI based driver to controll a string of WS281x
// LEDs. There are better solutions for multiple (8, 16) strings.
package wsspi

import (
	"led"
	"led/internal"
)

// Pixel represents SPI data that need to be send to the WS281x controller to
// set the color of one LED.
type Pixel struct {
	data [12]byte
}

// Bytes returns reference to the internal storage of p.
func (p *Pixel) Bytes() []byte {
	return p.data[:]
}

type ColorOrder byte

const (
	RGB ColorOrder = iota
	GRB
)

const zero = 0x88

func (co ColorOrder) pixel(r, g, b byte) Pixel {
	var p Pixel
	switch co {
	case RGB:
		for n := uint(0); n < 4; n++ {
			p.data[3-n] = zero | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4
			p.data[7-n] = zero | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4
			p.data[11-n] = zero | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4
		}
	case GRB:
		for n := uint(0); n < 4; n++ {
			p.data[3-n] = byte(zero | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4)
			p.data[7-n] = byte(zero | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4)
			p.data[11-n] = byte(zero | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4)
		}
	}
	return p
}

// PixelRaw returns a pixel with color set to c without gamma correction.
func (co ColorOrder) RawPixel(c led.Color) Pixel {
	return co.pixel(c.Red(), c.Green(), c.Blue())
}

// Pixel returns a pixel with color set to c with gamma correction.
func (co ColorOrder) Pixel(c led.Color) Pixel {
	r := internal.Gamma8(c.Red())
	g := internal.Gamma8(c.Green())
	b := internal.Gamma8(c.Blue())
	return co.pixel(r, g, b)
}
