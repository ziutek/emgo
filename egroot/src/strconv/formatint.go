package strconv

import (
	"io"
)

func digit(d, toA uint32) byte {
	if d < 10 {
		return byte(d + '0')
	}
	return byte(d + toA)
}

func fixBase(base int) (u, toA uint32) {
	toA = 'a' - 10
	if base < 0 {
		base = -base
		toA = 'A' - 10
	}
	if base < 2 || base > 36 {
		panicBase()
	}
	return uint32(base), toA
}

func formatUint32(buf []byte, u uint32, base int) int {
	b32, toA := fixBase(base)
	n := len(buf)
	for u != 0 {
		if n == 0 {
			panicBuffer()
		}
		n--
		newU := u / b32
		buf[n] = digit(u%b32, toA)
		u = newU
	}
	if n == len(buf) {
		n--
		buf[n] = '0'
	}
	return n
}

func formatUint64(buf []byte, u uint64, base int) int {
	b32, toA := fixBase(base)
	n := len(buf)
	for u != 0 {
		if n == 0 {
			panicBuffer()
		}
		n--
		newU := u / uint64(b32)
		buf[n] = digit(uint32(u%uint64(b32)), toA)
		u = newU
	}
	if n == len(buf) {
		n--
		buf[n] = '0'
	}
	return n
}

func WriteUint32(w io.Writer, u uint32, base, width int, pad rune) (int, error) {
	var buf [32]byte
	n := formatUint32(buf[:], u, base)
	return writePadded(w, buf[n:], width, pad)
}

// WriteInt32 works like WriteInt but is optimized for 32-bit numbers.
func WriteInt32(w io.Writer, i int32, base, width int, pad rune) (int, error) {
	var (
		buf [33]byte
		n   int
	)
	if i >= 0 {
		n = formatUint32(buf[:], uint32(i), base)
	} else {
		n = formatUint32(buf[:], uint32(-i), base) - 1
		buf[n] = '-'
	}
	return writePadded(w, buf[n:], width, pad)
}

func WriteUint64(w io.Writer, u uint64, base, width int, pad rune) (int, error) {
	var buf [64]byte
	n := formatUint64(buf[:], u, base)
	return writePadded(w, buf[n:], width, pad)
}

// WriteInt64 works like WriteInt.
func WriteInt64(w io.Writer, i int64, base, width int, pad rune) (int, error) {
	var (
		buf [65]byte
		n   int
	)
	if i >= 0 {
		n = formatUint64(buf[:], uint64(i), base)
	} else {
		n = formatUint64(buf[:], uint64(-i), base) - 1
		buf[n] = '-'
	}
	return writePadded(w, buf[n:], width, pad)
}

// WriteInt writes text representation of i to w using 2 <= base <= 36.
// If width > 0 then written value is right-justified, otherwise it is
// left-justified.
func WriteInt(w io.Writer, i, base, width int, pad rune) (int, error) {
	if intSize <= 4 {
		return WriteInt32(w, int32(i), base, width, pad)
	} else {
		return WriteInt64(w, int64(i), base, width, pad)
	}
}

func WriteUint(w io.Writer, u uint, base, width int, pad rune) (int, error) {
	if intSize <= 4 {
		return WriteUint32(w, uint32(u), base, width, pad)
	} else {
		return WriteUint64(w, uint64(u), base, width, pad)
	}
}

func WriteUintptr(w io.Writer, u uintptr, base, width int, pad rune) (int, error) {
	if ptrSize <= 4 {
		return WriteUint32(w, uint32(u), base, width, pad)
	} else {
		return WriteUint64(w, uint64(u), base, width, pad)
	}
}
