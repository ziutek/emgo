package ws281x

import (
	"bytes"
)

//  can be used to implement a frame buffer for SPI based WS281x driver.
// SPIFB encoding uses one byte of memory to encode two WS281x bits
// (12 bytes/pixel).
type SPIFB struct {
	data []byte
}

// MakeSPIFB allocates memory for string of n pixels.
func MakeSPIFB(n int) SPIFB {
	return SPIFB{make([]byte, n*12+1)}
}

// AsSPIFB returnfb SPIFB using s as data storage. For n pixels n*12+1 bytes are need.
func AsSPIFB(s []byte) SPIFB {
	s[len(s)-1] = 0 // STM32 SPI leaves MOSI high if last byte has LSBit set.
	return SPIFB{s}
}

// PixelSize returns pixel size in bytes (always 12).
func (_ SPIFB) PixelSize() int {
	return 12
}

// Len returnfb SPIFB length as number of pixels.
func (fb SPIFB) Len() int {
	return len(fb.data) / 12
}

// At returns slice of p that contains p.Len()-n pixels starting from n.
func (fb SPIFB) At(n int) SPIFB {
	return SPIFB{fb.data[n*12:]}
}

// Head returns slice of p that contains n pixels starting from 0.
func (fb SPIFB) Head(n int) SPIFB {
	return SPIFB{fb.data[:n*12]}
}

const zeroSPI = 0x88

// EncodeRGB encodes c to one pixel at begining of s in WS2811 RGB order.
func (fb SPIFB) EncodeRGB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	for n := uint(0); n < 4; n++ {
		fb.data[3-n] = byte(zeroSPI | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4)
		fb.data[7-n] = byte(zeroSPI | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4)
		fb.data[11-n] = byte(zeroSPI | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4)
	}
}

// EncodeGRB encodes c to one pixel at begining of s in WS2812 GRB order.
func (fb SPIFB) EncodeGRB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	for n := uint(0); n < 4; n++ {
		fb.data[3-n] = byte(zeroSPI | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4)
		fb.data[7-n] = byte(zeroSPI | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4)
		fb.data[11-n] = byte(zeroSPI | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4)
	}
}

// Bytes returns p's internal storage.
func (fb SPIFB) Bytes() []byte {
	return fb.data
}

// Write writes src to beginning of s.
func (fb SPIFB) Write(src SPIFB) {
	copy(fb.Bytes(), src.Bytes())
}

// Fill fills whole s using pattern p.
func (fb SPIFB) Fill(p SPIFB) {
	dst := fb.Bytes()
	src := p.Bytes()
	for i := 0; i < len(dst); i += copy(dst[i:], src) {
	}
}

// Clear clears whole s to black color.
func (fb SPIFB) Clear() {
	bytes.Fill(fb.data[:len(fb.data)-1], zeroSPI)
}
