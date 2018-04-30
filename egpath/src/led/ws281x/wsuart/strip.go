package wsuart

import (
	"bytes"
	"unsafe"
)

// Strip represents string of LEDs.
type Strip []Pixel

// Make returns cleared Strip of n pixels.
func Make(n int) Strip {
	s := make(Strip, n)
	s.Clear()
	return s
}

// Bytes returns reference to the internal storage of s.
func (s Strip) Bytes() []byte {
	return (*[1<<31 - 1]byte)(unsafe.Pointer(&s[0]))[:len(s)*8]
}

// Fill fills whole s with pixel p.
func (s Strip) Fill(p Pixel) {
	for i := range s {
		s[i] = p
	}
}

// Clear clears whole s to black color. It is faster than Fill(black).
func (s Strip) Clear() {
	bytes.Fill(s.Bytes(), zero)
}
