package ws281x

import (
	"bytes"
)

// UARTFB can be used to implement a frame buffer for UART based WS281x driver.
// UARTFB encoding uses one byte of memory to encode three WS281x bits
// (8 bytes/pixel).
type UARTFB struct {
	data []byte
}

// MakeUARTFB allocates memory for string of n pixels.
func MakeUARTFB(n int) UARTFB {
	return UARTFB{make([]byte, n*8)}
}

// AsUARTFB returns UARTFB using b as data storage. For n pixels n*8 bytes are need.
func AsUARTFB(s []byte) UARTFB {
	return UARTFB{s}
}

// PixelSize returns pixel size in bytes (always 8).
func (_ UARTFB) PixelSize() int {
	return 8
}

// Len returnfb UARTFB length as number of pixels.
func (fb UARTFB) Len() int {
	return len(fb.data) / 8
}

// At returns slice of p that contains p.Len()-n pixels starting from n.
func (fb UARTFB) At(n int) UARTFB {
	return UARTFB{fb.data[n*8:]}
}

// Head returns slice of p that contains n pixels starting from 0.
func (fb UARTFB) Head(n int) UARTFB {
	return UARTFB{fb.data[:n*8]}
}

const zeroUART = (6>>1 + 6<<2 + 6<<5)

// EncodeRGB encodes c to one pixel at begining of s in WS2811 RGB order.
func (fb UARTFB) EncodeRGB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	fb.data[0] = byte(zeroUART &^ (r>>7&1 | r>>3&8 | r<<1&0x40))
	fb.data[1] = byte(zeroUART &^ (r>>4&1 | r>>0&8 | r<<4&0x40))
	fb.data[2] = byte(zeroUART &^ (r>>1&1 | r<<3&8 | g>>1&0x40))
	fb.data[3] = byte(zeroUART &^ (g>>6&1 | g>>2&8 | g<<2&0x40))
	fb.data[4] = byte(zeroUART &^ (g>>3&1 | g<<1&8 | g<<5&0x40))
	fb.data[5] = byte(zeroUART &^ (g>>0&1 | b>>4&8 | b>>0&0x40))
	fb.data[6] = byte(zeroUART &^ (b>>5&1 | b>>1&8 | b<<3&0x40))
	fb.data[7] = byte(zeroUART &^ (b>>2&1 | b<<2&8 | b<<6&0x40))
}

// EncodeGRB encodes c to one pixel at begining of s in WS2812 GRB order.
func (fb UARTFB) EncodeGRB(c Color) {
	r, g, b := c.Red(), c.Green(), c.Blue()
	fb.data[0] = byte(zeroUART &^ (g>>7&1 | g>>3&8 | g<<1&0x40))
	fb.data[1] = byte(zeroUART &^ (g>>4&1 | g>>0&8 | g<<4&0x40))
	fb.data[2] = byte(zeroUART &^ (g>>1&1 | g<<3&8 | r>>1&0x40))
	fb.data[3] = byte(zeroUART &^ (r>>6&1 | r>>2&8 | r<<2&0x40))
	fb.data[4] = byte(zeroUART &^ (r>>3&1 | r<<1&8 | r<<5&0x40))
	fb.data[5] = byte(zeroUART &^ (r>>0&1 | b>>4&8 | b>>0&0x40))
	fb.data[6] = byte(zeroUART &^ (b>>5&1 | b>>1&8 | b<<3&0x40))
	fb.data[7] = byte(zeroUART &^ (b>>2&1 | b<<2&8 | b<<6&0x40))
}

// Bytes returns p's internal storage.
func (fb UARTFB) Bytes() []byte {
	return fb.data
}

// Write writes src to beginning of s.
func (fb UARTFB) Write(src UARTFB) {
	copy(fb.Bytes(), src.Bytes())
}

// Fill fills whole s using pattern p.
func (fb UARTFB) Fill(p UARTFB) {
	dst := fb.Bytes()
	src := p.Bytes()
	for i := 0; i < len(dst); i += copy(dst[i:], src) {
	}
}

// Clear clears whole s to black color.
func (fb UARTFB) Clear() {
	bytes.Fill(fb.data, zeroUART)
}
