package ws281x

import (
	"bytes"
)

// FB2 can be used to implement a frame buffer for SPI based WFB281x driver. FB2
// encoding uses one byte of memory to encode two WFB281x bits (12 bytes/pixel).
type FB2 struct {
	data []byte
}

// MakeFB2 allocates memory for string of n pixels.
func MakeFB2(n int) FB2 {
	return FB2{make([]byte, n*12+1)}
}

// AsFB2 returns FB2 using b as data storage. For n pixels n*12+1 bytes are need.
func AsFB2(b []byte) FB2 {
	b[len(b)-1] = 0 // STM32 SPI leaves MOSI high if last byte has LSBit set.
	return FB2{b}
}

// PixelSize returns pixel size.
func (_ FB2) PixelSize() int {
	return 12
}

func (s FB2) Len() int {
	return len(s.data) / 12
}

// At returns slice of p that contains p.Len()-n pixels starting from n.
func (s FB2) At(n int) FB2 {
	return FB2{s.data[n*12:]}
}

// Head returns slice of p that contains n pixels starting from 0.
func (s FB2) Head(n int) FB2 {
	return FB2{s.data[:n*12]}
}

const zero2 = 0x88

// EncodeRGB encodes c to one pixel at begining of buf in WFB2811 RGB order.
func (s FB2) EncodeRGB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	for n := uint(0); n < 4; n++ {
		s.data[3-n] = byte(zero2 | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4)
		s.data[7-n] = byte(zero2 | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4)
		s.data[11-n] = byte(zero2 | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4)
	}
}

// EncodeGRB encodes c to one pixel at begining of buf in WFB2812 GRB order.
func (s FB2) EncodeGRB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	for n := uint(0); n < 4; n++ {
		s.data[3-n] = byte(zero2 | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4)
		s.data[7-n] = byte(zero2 | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4)
		s.data[11-n] = byte(zero2 | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4)
	}
}

// Bytes returns p's internal storage.
func (s FB2) Bytes() []byte {
	return s.data
}

// Write writes src to beginning of p.
func (s FB2) Write(src FB2) {
	copy(s.Bytes(), src.Bytes())
}

// Fill fills whole s using pattern p.
func (s FB2) Fill(p FB2) {
	sb := s.Bytes()
	pb := p.Bytes()
	for i := 0; i < len(sb); i += copy(sb[i:], pb) {
	}
}

// Clear clears whole s to black color.
func (s FB2) Clear() {
	bytes.Fill(s.data[:len(s.data)-1], zero2)
}
