package strconv

import (
	"io"
)

func WriteString(w io.Writer, s string, width int, zeros bool) (int, error) {
	right := width < 0
	if right {
		width = -width
	}
	extn := width - len(s)
	var (
		n   int
		err error
	)
	if extn > 0 && right {
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
	if extn > 0 && !right {
		m, err = padd(w, pspaces, extn)
		n += m
	}
	return n, err
}
