package strconv

import (
	"io"
)

func panicBuffer() {
	panic("strconv: buffer too short")
}

func panicBase() {
	panic("strconv: unsupported base")
}

func writeRuneN(w io.Writer, r rune, n int) (int, error) {
	var (
		m     int
		chars [8]byte
	)
	// BUG: Casting rune to byte.
	for i := range chars {
		chars[i] = byte(r)
	}
	for {
		if n <= len(chars) {
			k, err := w.Write(chars[:n])
			return m + k, err
		}
		k, err := w.Write(chars[:])
		m += k
		if err != nil {
			return m, err
		}
		n -= len(chars)
	}
}

func writePadded(w io.Writer, b []byte, width int, pad rune) (int, error) {
	left := width < 0
	if left {
		width = -width
	}
	extn := width - len(b)
	var (
		m, n int
		err  error
	)
	if extn > 0 && !left {
		if pad == '0' {
			if b[0] == '-' {
				n, err = w.Write(b[:1])
				if err != nil {
					return n, err
				}
				b = b[1:]
			}
		}
		m, err = writeRuneN(w, pad, extn)
		n += m
		if err != nil {
			return n, err
		}
	}
	if len(b) > 0 {
		m, err = w.Write(b)
		n += m
		if err != nil {
			return n, err
		}
	}
	if extn > 0 && left {
		m, err = writeRuneN(w, ' ', extn)
		n += m
	}
	return n, err
}
