package ws281x

import (
	"bytes"
)

// FBU can be used to implement a frame buffer for UART based WS281x driver. FBU
// encoding uses one byte of memory to encode three WS281x bits (8 bytes/pixel).
type FBU struct {
	data []byte
}

// MakeFBU allocates memory for string of n pixels.
func MakeFBU(n int) FBU {
	return FBU{make([]byte, n*8)}
}

// AsFBU returns FBU using b as data storage. For n pixels n*8 bytes are need.
func AsFBU(b []byte) FBU {
	return FBU{b}
}

// PixelSize returns pixel size.
func (_ FBU) PixelSize() int {
	return 8
}

func (s FBU) Len() int {
	return len(s.data) / 8
}

// At returns slice of p that contains p.Len()-n pixels starting from n.
func (s FBU) At(n int) FBU {
	return FBU{s.data[n*8:]}
}

// Head returns slice of p that contains n pixels starting from 0.
func (s FBU) Head(n int) FBU {
	return FBU{s.data[:n*8]}
}

const zeroUART = (6>>1 + 6<<2 + 6<<5)

// EncodeRGB encodes c to one pixel at begining of buf in WS2811 RGB order.
func (s FBU) EncodeRGB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	s.data[0] = byte(zeroUART &^ (r>>7&1 | r>>3&8 | r<<1&0x40))
	s.data[1] = byte(zeroUART &^ (r>>4&1 | r>>0&8 | r<<4&0x40))
	s.data[2] = byte(zeroUART &^ (r>>1&1 | r<<3&8 | g>>1&0x40))
	s.data[3] = byte(zeroUART &^ (g>>6&1 | g>>2&8 | g<<2&0x40))
	s.data[4] = byte(zeroUART &^ (g>>3&1 | g<<1&8 | g<<5&0x40))
	s.data[5] = byte(zeroUART &^ (g>>0&1 | b>>4&8 | b>>0&0x40))
	s.data[6] = byte(zeroUART &^ (b>>5&1 | b>>1&8 | b<<3&0x40))
	s.data[7] = byte(zeroUART &^ (b>>2&1 | b<<2&8 | b<<6&0x40))
}

// EncodeGRB encodes c to one pixel at begining of buf in WS2812 GRB order.
func (s FBU) EncodeGRB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	s.data[0] = byte(zeroUART &^ (g>>7&1 | g>>3&8 | g<<1&0x40))
	s.data[1] = byte(zeroUART &^ (g>>4&1 | g>>0&8 | g<<4&0x40))
	s.data[2] = byte(zeroUART &^ (g>>1&1 | g<<3&8 | r>>1&0x40))
	s.data[3] = byte(zeroUART &^ (r>>6&1 | r>>2&8 | r<<2&0x40))
	s.data[4] = byte(zeroUART &^ (r>>3&1 | r<<1&8 | r<<5&0x40))
	s.data[5] = byte(zeroUART &^ (r>>0&1 | b>>4&8 | b>>0&0x40))
	s.data[6] = byte(zeroUART &^ (b>>5&1 | b>>1&8 | b<<3&0x40))
	s.data[7] = byte(zeroUART &^ (b>>2&1 | b<<2&8 | b<<6&0x40))
}

// Bytes returns p's internal storage.
func (s FBU) Bytes() []byte {
	return s.data
}

// Write writes src to beginning of p.
func (s FBU) Write(src FBU) {
	copy(s.Bytes(), src.Bytes())
}

// Fill fills whole s using pattern p.
func (s FBU) Fill(p FBU) {
	sb := s.Bytes()
	pb := p.Bytes()
	for i := 0; i < len(sb); i += copy(sb[i:], pb) {
	}
}

// Clear clears whole s to black color.
func (s FBU) Clear() {
	bytes.Fill(s.data, zeroUART)
}
