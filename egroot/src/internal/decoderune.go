package internal

// Code borrowed from Go standard library unicode/utf8 package.
// Original copyrigth:
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

const (
	RuneError = '\uFFFD'     // the "error" Rune or "Unicode replacement character"
	MaxRune   = '\U0010FFFF' // Maximum valid Unicode code point.
)

// Code points in the surrogate range are not valid for UTF-8.
const (
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
)

const (
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000

	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1
)

func DecodeRune(s string) (r rune, size int, short bool) {
	n := len(s)
	if n < 1 {
		return RuneError, 0, true
	}
	c0 := s[0]

	// 1-byte, 7-bit sequence?
	if c0 < tx {
		return rune(c0), 1, false
	}

	// unexpected continuation byte?
	if c0 < t2 {
		return RuneError, 1, false
	}

	// need first continuation byte
	if n < 2 {
		return RuneError, 1, true
	}
	c1 := s[1]
	if c1 < tx || t2 <= c1 {
		return RuneError, 1, false
	}

	// 2-byte, 11-bit sequence?
	if c0 < t3 {
		r = rune(c0&mask2)<<6 | rune(c1&maskx)
		if r <= rune1Max {
			return RuneError, 1, false
		}
		return r, 2, false
	}

	// need second continuation byte
	if n < 3 {
		return RuneError, 1, true
	}
	c2 := s[2]
	if c2 < tx || t2 <= c2 {
		return RuneError, 1, false
	}

	// 3-byte, 16-bit sequence?
	if c0 < t4 {
		r = rune(c0&mask3)<<12 | rune(c1&maskx)<<6 | rune(c2&maskx)
		if r <= rune2Max {
			return RuneError, 1, false
		}
		if surrogateMin <= r && r <= surrogateMax {
			return RuneError, 1, false
		}
		return r, 3, false
	}

	// need third continuation byte
	if n < 4 {
		return RuneError, 1, true
	}
	c3 := s[3]
	if c3 < tx || t2 <= c3 {
		return RuneError, 1, false
	}

	// 4-byte, 21-bit sequence?
	if c0 < t5 {
		r = rune(c0&mask4)<<18 | rune(c1&maskx)<<12 | rune(c2&maskx)<<6 | rune(c3&maskx)
		if r <= rune3Max || MaxRune < r {
			return RuneError, 1, false
		}
		return r, 4, false
	}

	// error
	return RuneError, 1, false
}
