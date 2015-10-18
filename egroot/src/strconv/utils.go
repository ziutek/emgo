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

const (
	pspaces = "        "
	pzeros  = "00000000"
)

func padd(w io.Writer, chars string, n int) (int, error) {
	var m int
	for {
		if n <= len(chars) {
			k, err := io.WriteString(w, chars[:n])
			return m + k, err
		}
		k, err := io.WriteString(w, chars)
		m += k
		if err != nil {
			return m, err
		}
		n -= len(chars)
	}
}

func writePadded(w io.Writer, b []byte, width int, zeros bool) (int, error) {
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
		if zeros {
			if b[0] == '-' {
				n, err = w.Write(b[:1])
				if err != nil {
					return n, err
				}
				b = b[1:]
			}
			m, err = padd(w, pzeros, extn)
		} else {
			m, err = padd(w, pspaces, extn)
		}
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
		m, err = padd(w, pspaces, extn)
		n += m
	}
	return n, err
}
