package ws281x

import (
	"bytes"
)

// FB3 can be used to implement a frame buffer for UART based WS281x driver. FB3
// encoding uses one byte of memory to encode three WS281x bits (8 bytes/pixel).
type FB3 struct {
	data []byte
}

// MakeFB3 allocates memory for string of n pixels.
func MakeFB3(n int) FB3 {
	return FB3{make([]byte, n*8)}
}

// AsFB3 returns FB3 using b as data storage. For n pixels n*8 bytes are need.
func AsFB3(b []byte) FB3 {
	return FB3{b}
}

// PixelSize returns pixel size.
func (_ FB3) PixelSize() int {
	return 8
}

func (s FB3) Len() int {
	return len(s.data) / 8
}

// At returns slice of p that contains p.Len()-n pixels starting from n.
func (s FB3) At(n int) FB3 {
	return FB3{s.data[n*8:]}
}

// Head returns slice of p that contains n pixels starting from 0.
func (s FB3) Head(n int) FB3 {
	return FB3{s.data[:n*8]}
}

const zero3 = (6>>1 + 6<<2 + 6<<5)

// EncodeRGB encodes c to one pixel at begining of buf in WS2811 RGB order.
func (s FB3) EncodeRGB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	s.data[0] = byte(zero3 &^ (r>>7&1 | r>>3&8 | r<<1&0x40))
	s.data[1] = byte(zero3 &^ (r>>4&1 | r>>0&8 | r<<4&0x40))
	s.data[2] = byte(zero3 &^ (r>>1&1 | r<<3&8 | g>>1&0x40))
	s.data[3] = byte(zero3 &^ (g>>6&1 | g>>2&8 | g<<2&0x40))
	s.data[4] = byte(zero3 &^ (g>>3&1 | g<<1&8 | g<<5&0x40))
	s.data[5] = byte(zero3 &^ (g>>0&1 | b>>4&8 | b>>0&0x40))
	s.data[6] = byte(zero3 &^ (b>>5&1 | b>>1&8 | b<<3&0x40))
	s.data[7] = byte(zero3 &^ (b>>2&1 | b<<2&8 | b<<6&0x40))
}

// EncodeGRB encodes c to one pixel at begining of buf in WS2812 GRB order.
func (s FB3) EncodeGRB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	s.data[0] = byte(zero3 &^ (g>>7&1 | g>>3&8 | g<<1&0x40))
	s.data[1] = byte(zero3 &^ (g>>4&1 | g>>0&8 | g<<4&0x40))
	s.data[2] = byte(zero3 &^ (g>>1&1 | g<<3&8 | r>>1&0x40))
	s.data[3] = byte(zero3 &^ (r>>6&1 | r>>2&8 | r<<2&0x40))
	s.data[4] = byte(zero3 &^ (r>>3&1 | r<<1&8 | r<<5&0x40))
	s.data[5] = byte(zero3 &^ (r>>0&1 | b>>4&8 | b>>0&0x40))
	s.data[6] = byte(zero3 &^ (b>>5&1 | b>>1&8 | b<<3&0x40))
	s.data[7] = byte(zero3 &^ (b>>2&1 | b<<2&8 | b<<6&0x40))
}

// Bytes returns p's internal storage.
func (s FB3) Bytes() []byte {
	return s.data
}

// Write writes src to beginning of p.
func (s FB3) Write(src FB3) {
	copy(s.Bytes(), src.Bytes())
}

// Fill fills whole s using pattern p.
func (s FB3) Fill(p FB3) {
	sb := s.Bytes()
	pb := p.Bytes()
	for i := 0; i < len(sb); i += copy(sb[i:], pb) {
	}
}

// Clear clears whole s to black color.
func (s FB3) Clear() {
	bytes.Fill(s.data, zero3)
}
