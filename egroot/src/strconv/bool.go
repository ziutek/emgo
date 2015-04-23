package strconv

import "bytes"

func panicBuffer() {
	panic("strconv: buffer too short")
}

func panicBase() {
	panic("strconv: unsupported base")
}

func Btoa(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// FormatBool stores text representation of b in buf using format specified by
// base:
//	|base| == 1: 0 / 1,
//  |base| == 2: fales / true.
// Unused portion of the buffer is filed with spaces.
// If base > 0 then formatted value is left-justified and FormatBool returns
// its length. If base < 0 then formatted value is right-justified and
// FormatBool returns offset to its first char.
func FormatBool(buf []byte, b bool, base int) int {
	right := base < 0
	if right {
		base = -base
	}
	if base != 1 && base != 2 {
		panicBase()
	}
	if len(buf) < 4*base-3 {
		panicBuffer()
	}
	var str string
	if m, n := base-1, 4*(base-2)+5; b {
		str = "1true"[m:n]
	} else {
		str = "0false"[m:][:n]
	}
	if !right {
		copy(buf, str)
		bytes.Fill(buf[:len(str)], ' ')
		return len(str)
	}
	n := len(buf) - len(str)
	copy(buf[n:], str)
	bytes.Fill(buf[:n], ' ')
	return n
}
