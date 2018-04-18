// Package wsuart allows to use UART based driver to controll a string of WS281x
// LEDs. There are better solutions for multiple (8, 16) strings.
package wsuart

import "led"

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

const zero = (6>>1 + 6<<2 + 6<<5)

// RawPixel returns a pixel with color set to raw (r, g, b).
func (c ColorOrder) RawPixel(r, g, b byte) Pixel {
	switch c {
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

// Pixel returns a pixel with color set to (r, g, b) with gamma corection.
func (c ColorOrder) Pixel(r, g, b byte) Pixel {
	return c.RawPixel(led.Gamma8(r), led.Gamma8(g), led.Gamma8(b))
}
