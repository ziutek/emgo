package onewire

import (
	"fmt"
	"io"
)

type Dev [8]byte

func (d Dev) Format(w io.Writer, p ...int) (n int, err error) {
	sep := fmt.Rune('-')
	for i, b := range d {
		var m int
		if i != 0 {
			m, err = sep.Format(w)
			n += m
			if err != nil {
				return
			}
		}
		m, err = fmt.Byte(b).Format(w, 16)
		n += m
		if err != nil {
			return
		}
	}
	return
}
