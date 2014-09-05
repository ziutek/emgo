package fmt

import (
	"io"
	"unsafe"
	"strconv"
)

// Formatter is the interface implemented by values hat can be printed.
type Formatter interface {
	// Format writes text representation of value to w. It can use optional
	// parameters p to configure formatter (base, width, precision, etc...).
	// Format returns numbe of bytes written and any write error encountered.
	Format(w io.Writer, p ...int) (n int, err error)
}

type Int32

func (i Int32) Format(w io.Writer, a ...int) (int, error) {
	base := 10
	if len(a) > 0 {
		base = a[0]
	}
	width := 0
	if len(a) > 1 {
		width = a[1]
	}
	var buf []byte
	if width < 10 {
		buf = stack.Bytes(10)
	} else {
		buf = stack.Bytes(width)
	}
	first := Utoa32(buf, u, base)
	return w.Write(buf[first:])

}