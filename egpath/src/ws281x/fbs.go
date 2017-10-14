package ws281x

import (
	"bytes"
)

// FBS can be used to implement a frame buffer for SPI based WS281x driver. FBS
// encoding uses one byte of memory to encode two WS281x bits (12 bytes/pixel).
type FBS struct {
	data []byte
}

// MakeFBS allocates memory for string of n pixels.
func MakeFBS(n int) FBS {
	return FBS{make([]byte, n*12+1)}
}

// AsFBS returns FBS using b as data storage. For n pixels n*12+1 bytes are need.
func AsFBS(b []byte) FBS {
	b[len(b)-1] = 0 // STM32 SPI leaves MOSI high if last byte has LSBit set.
	return FBS{b}
}

// PixelSize returns pixel size in bytes (always 12).
func (_ FBS) PixelSize() int {
	return 12
}

// Len returns FBS length as number of pixels.
func (s FBS) Len() int {
	return len(s.data) / 12
}

// At returns slice of p that contains p.Len()-n pixels starting from n.
func (s FBS) At(n int) FBS {
	return FBS{s.data[n*12:]}
}

// Head returns slice of p that contains n pixels starting from 0.
func (s FBS) Head(n int) FBS {
	return FBS{s.data[:n*12]}
}

const zeroSPI = 0x88

// EncodeRGB encodes c to one pixel at begining of s in WS2811 RGB order.
func (s FBS) EncodeRGB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	for n := uint(0); n < 4; n++ {
		s.data[3-n] = byte(zeroSPI | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4)
		s.data[7-n] = byte(zeroSPI | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4)
		s.data[11-n] = byte(zeroSPI | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4)
	}
}

// EncodeGRB encodes c to one pixel at begining of s in WS2812 GRB order.
func (s FBS) EncodeGRB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	for n := uint(0); n < 4; n++ {
		s.data[3-n] = byte(zeroSPI | g>>(2*n+1)&1<<6 | g>>(2*n)&1<<4)
		s.data[7-n] = byte(zeroSPI | r>>(2*n+1)&1<<6 | r>>(2*n)&1<<4)
		s.data[11-n] = byte(zeroSPI | b>>(2*n+1)&1<<6 | b>>(2*n)&1<<4)
	}
}

// Bytes returns p's internal storage.
func (s FBS) Bytes() []byte {
	return s.data
}

// Write writes src to beginning of s.
func (s FBS) Write(src FBS) {
	copy(s.Bytes(), src.Bytes())
}

// Fill fills whole s using pattern p.
func (s FBS) Fill(p FBS) {
	sb := s.Bytes()
	pb := p.Bytes()
	for i := 0; i < len(sb); i += copy(sb[i:], pb) {
	}
}

// Clear clears whole s to black color.
func (s FBS) Clear() {
	bytes.Fill(s.data[:len(s.data)-1], zeroSPI)
}
