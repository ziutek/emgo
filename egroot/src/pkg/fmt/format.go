package fmt

import (
	"io"
	"stack"
	"strconv"
)

// Formatter is the interface implemented by values that can present itself in text
// form.
type Formatter interface {
	// Format writes text representation of value to w. It can use optional
	// parameters p to configure formatter (base, width, precision, etc...).
	// Format returns numbe of bytes written and any write error encountered.
	Format(w io.Writer, p ...int) (n int, err error)
}

type Uint32 uint

func (u Uint32) Format(w io.Writer, a ...int) (int, error) {
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
	first := strconv.Utoa32(buf, uint32(u), base)
	if f := len(buf) - width; first > f {
		first = f
	}
	return w.Write(buf[first:])
}

type Int32 int32

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
	if width < 11 {
		buf = stack.Bytes(11)
	} else {
		buf = stack.Bytes(width)
	}
	first := strconv.Itoa32(buf, int32(i), base)
	if f := len(buf) - width; first > f {
		first = f
	}
	return w.Write(buf[first:])
}

type Byte byte

func (b Byte) Format(w io.Writer, a ...int) (int, error) {
	return Uint32(b).Format(w, a...)
}

type Int8 int8

func (i Int8) Format(w io.Writer, a ...int) (int, error) {
	return Int32(i).Format(w, a...)
}
