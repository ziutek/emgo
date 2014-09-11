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

type Uint64 uint64

func (u Uint64) Format(w io.Writer, a ...int) (int, error) {
	base := 10
	if len(a) > 0 {
		base = a[0]
	}
	width := 0
	if len(a) > 1 {
		width = a[1]
	}
	var buf []byte
	if width < 20 {
		buf = stack.Bytes(20)
	} else {
		buf = stack.Bytes(width)
	}
	first := strconv.Utoa64(buf, uint64(u), base)
	if f := len(buf) - width; first > f {
		first = f
	}
	return w.Write(buf[first:])
}

type Int64 int64

func (u Int64) Format(w io.Writer, a ...int) (int, error) {
	base := 10
	if len(a) > 0 {
		base = a[0]
	}
	width := 0
	if len(a) > 1 {
		width = a[1]
	}
	var buf []byte
	if width < 21 {
		buf = stack.Bytes(21)
	} else {
		buf = stack.Bytes(width)
	}
	first := strconv.Itoa64(buf, int64(u), base)
	if f := len(buf) - width; first > f {
		first = f
	}
	return w.Write(buf[first:])
}

type Uint32 uint32

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

type Uint16 uint16

func (i Uint16) Format(w io.Writer, a ...int) (int, error) {
	return Uint32(i).Format(w, a...)
}

type Int16 int16

func (i Int16) Format(w io.Writer, a ...int) (int, error) {
	return Int32(i).Format(w, a...)
}

type Byte byte

func (b Byte) Format(w io.Writer, a ...int) (int, error) {
	return Uint32(b).Format(w, a...)
}

type Int8 int8

func (i Int8) Format(w io.Writer, a ...int) (int, error) {
	return Int32(i).Format(w, a...)
}
