package strconv

import (
	"io"
)

func WriteString(w io.Writer, s string, width int, zeros bool) (int, error) {
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
		if zeros {
			n, err = padd(w, pzeros, extn)
		} else {
			n, err = padd(w, pspaces, extn)
		}
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
		m, err = padd(w, pspaces, extn)
		n += m
	}
	return n, err
}

func WriteBytes(w io.Writer, s []byte, width int, zeros bool) (int, error) {
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
		if zeros {
			n, err = padd(w, pzeros, extn)
		} else {
			n, err = padd(w, pspaces, extn)
		}
		if err != nil {
			return n, err
		}
	}
	m, err := w.Write(s)
	n += m
	if err != nil {
		return n, err
	}
	if extn > 0 && left {
		m, err = padd(w, pspaces, extn)
		n += m
	}
	return n, err
}
