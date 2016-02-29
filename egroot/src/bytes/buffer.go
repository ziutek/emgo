package bytes

import (
	"errors"
)

var ErrTooLarge = errors.New("bytes.Buffer: too large")

type Buffer struct {
	buf   []byte
	off   int
	fixed bool
}

func NewBuffer(buf []byte) *Buffer {
	b := new(Buffer)
	b.buf = buf
	return b
}

func MakeBuffer(buf []byte, fixed bool) Buffer {
	return Buffer{buf: buf, fixed: fixed}
}

func (b *Buffer) Len() int { return len(b.buf) - b.off }

func (b *Buffer) Cap() int { return cap(b.buf) }

func (b *Buffer) Bytes() []byte { return b.buf[b.off:] }

func (b *Buffer) Truncate(n int) {
	if n == 0 {
		b.off = 0
	}
	b.buf = b.buf[0 : b.off+n]
}

func (b *Buffer) Reset() { b.Truncate(0) }

// Grow tries to grow the buffer to make space for n more bytes. In case of
// fixed buffer it can not allocate more space than current buffer has. Grow
// returns index wher bytes should be copied.
func grow(b *Buffer, n int) int {
	m := len(b.buf)
	n += m
	if n > cap(b.buf) {
		if !b.fixed {
			buf := make([]byte, n)
			copy(buf, b.buf)
			b.buf = buf
			return m
		}
		n = cap(b.buf)
	}
	b.buf = b.buf[:n]
	return m
}

func (b *Buffer) Write(s []byte) (int, error) {
	m := grow(b, len(s))
	n := copy(b.buf[m:], s)
	if n != len(s) {
		return n, ErrTooLarge
	}
	return n, nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	m := grow(b, len(s))
	n := copy(b.buf[m:], s)
	if n != len(s) {
		return n, ErrTooLarge
	}
	return n, nil
}
