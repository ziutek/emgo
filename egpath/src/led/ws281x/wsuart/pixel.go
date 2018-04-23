// Package wsuart allows to use UART based driver to controll a string of WS281x
// LEDs. There are better solutions for multiple (8, 16) strings.
package wsuart

import (
	"led"
	"led/internal"
)

// Pixel represents UART data that need to be send to the WS281x controller to
// set the color of one LED.
type Pixel struct {
	data [8]byte
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

const zero = 6>>1 | 6<<2 | 6<<5

func (co ColorOrder) pixel(r, g, b byte) Pixel {
	switch co {
	case RGB:
		return Pixel{[8]byte{
			zero &^ (r>>7&1 | r>>3&8 | r<<1&0x40),
			zero &^ (r>>4&1 | r>>0&8 | r<<4&0x40),
			zero &^ (r>>1&1 | r<<3&8 | g>>1&0x40),
			zero &^ (g>>6&1 | g>>2&8 | g<<2&0x40),
			zero &^ (g>>3&1 | g<<1&8 | g<<5&0x40),
			zero &^ (g>>0&1 | b>>4&8 | b>>0&0x40),
			zero &^ (b>>5&1 | b>>1&8 | b<<3&0x40),
			zero &^ (b>>2&1 | b<<2&8 | b<<6&0x40),
		}}
	case GRB:
		return Pixel{[8]byte{
			zero &^ (g>>7&1 | g>>3&8 | g<<1&0x40),
			zero &^ (g>>4&1 | g>>0&8 | g<<4&0x40),
			zero &^ (g>>1&1 | g<<3&8 | r>>1&0x40),
			zero &^ (r>>6&1 | r>>2&8 | r<<2&0x40),
			zero &^ (r>>3&1 | r<<1&8 | r<<5&0x40),
			zero &^ (r>>0&1 | b>>4&8 | b>>0&0x40),
			zero &^ (b>>5&1 | b>>1&8 | b<<3&0x40),
			zero &^ (b>>2&1 | b<<2&8 | b<<6&0x40),
		}}
	}
	return Pixel{}
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
