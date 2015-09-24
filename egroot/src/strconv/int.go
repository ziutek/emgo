package strconv

import (
	"io"
	"unsafe"
)

const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

func fixBase(base int) uint32 {
	if base < 0 {
		base = -base
	}
	if base < 2 || base > len(digits) {
		panicBase()
	}
	return uint32(base)
}

func formatUint32(buf []byte, u, base uint32) int {
	n := len(buf)
	for u != 0 {
		if n == 0 {
			panicBuffer()
		}
		n--
		newU := u / base
		buf[n] = digits[u-newU*base]
		u = newU
	}
	if n == len(buf) {
		n--
		buf[n] = '0'
	}
	return n
}

func formatUint64(buf []byte, u uint64, base uint32) int {
	n := len(buf)
	for u != 0 {
		if n == 0 {
			panicBuffer()
		}
		n--
		newU := u / uint64(base)
		buf[n] = digits[u-newU*uint64(base)]
		u = newU
	}
	if n == len(buf) {
		n--
		buf[n] = '0'
	}
	return n
}

func WriteUint32(w io.Writer, u uint32, base, width int) (int, error) {
	var buf [32]byte
	n := formatUint32(buf[:], u, fixBase(base))
	return writePadded(w, buf[n:], width, base < 0)
}

func WriteInt32(w io.Writer, i int32, base, width int) (int, error) {
	var (
		buf [33]byte
		n   int
	)
	b := fixBase(base)
	if i >= 0 {
		n = formatUint32(buf[:], uint32(i), b)
	} else {
		n = formatUint32(buf[:], uint32(-i), b) - 1
		buf[n] = '-'
	}
	return writePadded(w, buf[n:], width, base < 0)
}

func WriteUint64(w io.Writer, u uint64, base, width int) (int, error) {
	var buf [64]byte
	n := formatUint64(buf[:], u, fixBase(base))
	return writePadded(w, buf[n:], width, base < 0)
}

func WriteInt64(w io.Writer, i int64, base, width int) (int, error) {
	var (
		buf [65]byte
		n   int
	)
	b := fixBase(base)
	if i >= 0 {
		n = formatUint64(buf[:], uint64(i), b)
	} else {
		n = formatUint64(buf[:], uint64(-i), b) - 1
		buf[n] = '-'
	}
	return writePadded(w, buf[n:], width, base < 0)
}

// WriteInt writes text representation of i in buf using 2 <= |base| <= 36.
// If width > 0 then written value is left-justified, if width < 0 written
// value is right-justified. If base < 0 then right-justified value is prepended
// with zeros, instead it is padded with spaces.
func WriteInt(w io.Writer, i, base, width int) (int, error) {
	if unsafe.Sizeof(i) <= 4 {
		return WriteInt32(w, int32(i), base, width)
	} else {
		return WriteInt64(w, int64(i), base, width)
	}
}

func WriteUint(w io.Writer, u uint, base, width int) (int, error) {
	if unsafe.Sizeof(u) <= 4 {
		return WriteUint32(w, uint32(u), base, width)
	} else {
		return WriteUint64(w, uint64(u), base, width)
	}
}

func WriteUintptr(w io.Writer, u uintptr, base, width int) (int, error) {
	if unsafe.Sizeof(u) <= 4 {
		return WriteUint32(w, uint32(u), base, width)
	} else {
		return WriteUint64(w, uint64(u), base, width)
	}
}
