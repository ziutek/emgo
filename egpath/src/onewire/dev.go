package onewire

import (
	"fmt"
	"io"
)

type Dev uint64

func (d Dev) Format(w io.Writer, p ...int) (n int, err error) {
	crc := fmt.Byte(d >> 56)
	family := fmt.Byte(d)
	serial := fmt.Uint64(d>>8) & 0xffffffffffff
	sep := fmt.Rune('-')

	for _, f := range []fmt.Formatter{crc, sep, serial, sep, family} {
		var m int
		m, err = f.Format(w, 16)
		n += m
		if err != nil {
			return
		}
	}
	return
}
