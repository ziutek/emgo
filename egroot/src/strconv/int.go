package strconv

import "unsafe"

const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

func panicBuffer() {
	panic("strconv: buffer too short")
}

func panicBase() {
	panic("strconv: unsupported base")
}

func checkBase(base int) (int, bool) {
	move := false
	if base < 0 {
		base = -base
		move = true
	}
	if base < 2 || base > len(digits) {
		panicBase()
	}
	return base, move
}

func finish(buf []byte, n int, move bool) int {
	if n == len(buf) {
		if n == 0 {
			panicBuffer()
		}
		n--
		buf[n] = '0'
	}
	if n == 0 {
		return 0
	}
	i := 0
	end := n
	if move {
		i = copy(buf, buf[n:])
		end = len(buf)
		n = i
	}
	for i < end {
		buf[i] = ' '
		i++
	}
	return n
}

// FormatUint32 stores text representation of u in buf using 2 <= |base| <= 36.
// Unused portion of the buffer is filed with spaces. If base > 0 then formatted
// value is right-justified and FormatUint32 returns offset to its first char.
// If base < 0 then formatted value is left-justified and FormatUint32 returns
// its length.
func FormatUint32(buf []byte, u uint32, base int) int {
	base, move := checkBase(base)
	b := uint32(base)
	n := len(buf)
	for u != 0 {
		if n == 0 {
			panicBuffer()
		}
		n--
		newU := u / b
		buf[n] = digits[u-newU*b]
		u = newU
	}
	return finish(buf, n, move)
}

// FormatInt32 works like FormatUint32.
func FormatInt32(buf []byte, i int32, base int) int {
	if i >= 0 {
		return FormatUint32(buf, uint32(i), base)
	}
	if len(buf) < 2 {
		panicBuffer()
	}
	n := FormatUint32(buf[1:], uint32(-i), base)
	if base > 0 {
		buf[0] = ' '
		buf[n] = '-'
	} else {
		buf[0] = '-'
		n++
	}
	return n
}

// FormatUint64 works like FormatUint32
func FormatUint64(buf []byte, u uint64, base int) int {
	base, move := checkBase(base)
	b := uint64(base)
	n := len(buf)
	for u != 0 {
		if n == 0 {
			panicBuffer()
		}
		n--
		newU := u / b
		buf[n] = digits[u-newU*b]
		u = newU
	}
	return finish(buf, n, move)
}

// FormatInt64 works like FormatUint32.
func FormatInt64(buf []byte, i int64, base int) int {
	if i >= 0 {
		return FormatUint64(buf, uint64(i), base)
	}
	if len(buf) < 2 {
		panicBuffer()
	}
	n := FormatUint64(buf[1:], uint64(-i), base)
	if base > 0 {
		buf[0] = ' '
		buf[n] = '-'
	} else {
		buf[0] = '-'
		n++
	}
	return n
}

func FormatUint(buf []byte, u uint, base int) int {
	if unsafe.Sizeof(u) <= 4 {
		return FormatUint32(buf, uint32(u), base)
	}
	return FormatUint64(buf, uint64(u), base)
}

func FormatInt(buf []byte, i, base int) int {
	if unsafe.Sizeof(i) <= 4 {
		return FormatInt32(buf, int32(i), base)
	}
	return FormatInt64(buf, int64(i), base)
}
