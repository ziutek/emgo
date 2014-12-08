package fmt

import (
	"io"
	"stack"
	"strconv"
	"unsafe"
)

// Formatter is the interface implemented by values that can present itself in text
// form.
type Formatter interface {
	// Format writes text representation of value to w. It can use optional
	// parameters p to configure formatter (base, width, precision, etc...).
	// Format returns number of bytes written and any write error encountered.
	Format(w io.Writer, p ...int) (n int, err error)
}

func params(a []int, dbase int) (width int, alignr bool, base int) {
	if len(a) > 0 {
		width = a[0]
	}
	if alignr = (width < 0); alignr {
		width = -width
	}
	base = dbase
	if len(a) > 1 {
		base = a[1]
	}
	return
}
func writeNum(w io.Writer, buf []byte, first, width int, alignr bool) (n int, err error) {
	f := len(buf) - width
	if f > first {
		f = first
	}
	if alignr {
		return w.Write(buf[f:])
	}
	if n, err = w.Write(buf[first:]); err != nil {
		return
	}
	if m := first - f; m > 0 {
		spaces := stack.Bytes(m)
		for i := range spaces {
			spaces[i] = ' '
		}
		m, err = w.Write(spaces)
		n += m
	}
	return
}

type Uint64 uint64

func (u Uint64) Format(w io.Writer, a ...int) (int, error) {
	width, alignr, base := params(a, 10)
	var buf []byte
	if width < 20 {
		buf = stack.Bytes(20)
	} else {
		buf = stack.Bytes(width)
	}
	first := strconv.Utoa64(buf, uint64(u), base)
	return writeNum(w, buf, first, width, alignr)
}

type Int64 int64

func (i Int64) Format(w io.Writer, a ...int) (int, error) {
	width, alignr, base := params(a, 10)
	var buf []byte
	if width < 21 {
		buf = stack.Bytes(21)
	} else {
		buf = stack.Bytes(width)
	}
	first := strconv.Itoa64(buf, int64(i), base)
	return writeNum(w, buf, first, width, alignr)
}

type Uint32 uint32

func (u Uint32) Format(w io.Writer, a ...int) (int, error) {
	width, alignr, base := params(a, 10)
	var buf []byte
	if width < 10 {
		buf = stack.Bytes(10)
	} else {
		buf = stack.Bytes(width)
	}
	first := strconv.Utoa32(buf, uint32(u), base)
	return writeNum(w, buf, first, width, alignr)
}

type Int32 int32

func (i Int32) Format(w io.Writer, a ...int) (int, error) {
	width, alignr, base := params(a, 10)
	var buf []byte
	if width < 11 {
		buf = stack.Bytes(11)
	} else {
		buf = stack.Bytes(width)
	}
	first := strconv.Itoa32(buf, int32(i), base)
	return writeNum(w, buf, first, width, alignr)
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

type Uint int

func (u Uint) Format(w io.Writer, a ...int) (int, error) {
	if unsafe.Sizeof(u) <= 4 {
		return Uint32(u).Format(w, a...)
	}
	return Uint64(u).Format(w, a...)
}

type Uintptr uintptr

func (u Uintptr) Format(w io.Writer, a ...int) (int, error) {
	width, alignr, base := params(a, 16)
	var buf []byte
	need := 20
	if unsafe.Sizeof(u) <= 4 {
		need = 10
	}
	if width < need {
		buf = stack.Bytes(need)
	} else {
		buf = stack.Bytes(width)
	}
	var first int
	if unsafe.Sizeof(u) <= 4 {
		first = strconv.Utoa32(buf, uint32(u), base)
	} else {
		first = strconv.Utoa64(buf, uint64(u), base)
	}
	return writeNum(w, buf, first, width, alignr)
}

type Int int

func (i Int) Format(w io.Writer, a ...int) (int, error) {
	if unsafe.Sizeof(i) <= 4 {
		return Int32(i).Format(w, a...)
	}
	return Int64(i).Format(w, a...)
}

type Bytes []byte

func (b Bytes) Format(w io.Writer, a ...int) (n int, err error) {
	width := 0
	if len(a) > 0 {
		width = a[0]
	}
	alignr := false
	if alignr = width < 0; alignr {
		width = -width
	}
	var buf []byte
	if n := width - len(b); n > 0 {
		buf = stack.Bytes(n)
		for i := range buf {
			buf[i] = ' '
		}
	}
	var m int
	if alignr {
		n, err = w.Write(buf)
		if err != nil {
			return
		}
		m, err = w.Write(b)
	} else {
		n, err = w.Write(b)
		if err != nil {
			return
		}
		m, err = w.Write(buf)
	}
	n += m
	return
}

type Str string

func (s Str) Format(w io.Writer, a ...int) (int, error) {
	return Bytes(s).Format(w, a...)
}

type Rune rune

func (r Rune) Format(w io.Writer, a ...int) (int, error) {
	// BUG: need UTF8 support.
	return w.Write([]byte{byte(r)})
}

const (
	A Rune = '\a'
	B Rune = '\b'
	F Rune = '\f'
	N Rune = '\n'
	R Rune = '\r'
	S Rune = ' '
	T Rune = '\t'
	V Rune = '\v'
)

func Err(e error) Str {
	return Str(e.Error())
}
