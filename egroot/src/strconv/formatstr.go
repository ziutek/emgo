package strconv

import (
	"io"
)

func WriteString(w io.Writer, s string, width int, pad rune) (int, error) {
	left := width < 0
	if left {
		width = -width
	}
	extn := width - len(s)
	var (
		n   int
		err error
	)
	if extn > 0 && !left {
		n, err = writeRuneN(w, pad, extn)
		if err != nil {
			return n, err
		}
	}
	m, err := io.WriteString(w, s)
	n += m
	if err != nil {
		return n, err
	}
	if extn > 0 && left {
		if pad == '0' {
			pad = ' '
		}
		m, err = writeRuneN(w, pad, extn)
		n += m
	}
	return n, err
}
