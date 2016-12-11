package ws281x

import (
	"bytes"
)

// S2 can be used to implement a frame buffer for SPI based WS281x driver. S2
// encoding uses one byte of memory to encode two WS281x bits (12 bytes/pixel).
type S2 struct {
	data []byte
}

// MakeS2 allocates memory for string of n pixels.
func MakeS2(n int) S2 {
	return S2{make([]byte, n*12+1)}
}

// AsS2 returns S2 using b as data storage. For n pixels n*12+1 bytes are need.
func AsS2(b []byte) S2 {
	b[len(b)-1] = 0 // STM32 SPI leaves MOSI high if last byte has LSBit set.
	return S2{b}
}

// PixelSize returns pixel size.
func (_ S2) PixelSize() int {
	return 12
}

func (s S2) Len() int {
	return len(s.data) / 12
}

// At returns slice of p that contains p.Len()-n pixels starting from n.
func (s S2) At(n int) S2 {
	return S2{s.data[n*12:]}
}

// Head returns slice of p that contains n pixels starting from 0.
func (s S2) Head(n int) S2 {
	return S2{s.data[n*12:]}
}

const zero2 = 0x88

// EncodeRGB encodes c to one pixel at begining of buf in WS2811 RGB order.
func (s S2) EncodeRGB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	for n := uint(0); n < 4; n++ {
		s.data[3-n] = byte(zero2 | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4)
		s.data[7-n] = byte(zero2 | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4)
		s.data[11-n] = byte(zero2 | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4)
	}
}

// EncodeGRB encodes c to one pixel at begining of buf in WS2812 GRB order.
func (s S2) EncodeGRB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	for n := uint(0); n < 4; n++ {
		s.data[3-n] = byte(zero2 | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4)
		s.data[7-n] = byte(zero2 | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4)
		s.data[11-n] = byte(zero2 | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4)
	}
}

// Clear clears whole p to black color.
func (s S2) Clear() {
	bytes.Fill(s.data[:len(s.data)-1], zero2)
}

// Bytes returns p's internal storage.
func (s S2) Bytes() []byte {
	return s.data
}

// Write writes src to beginning of p.
func (s S2) Write(src S2) {
	copy(s.Bytes(), src.Bytes())
}
