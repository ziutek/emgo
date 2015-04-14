package strconv

import "unsafe"

const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

func fixBase(base int) int {
	if base < 0 {
		base = -base
	}
	if base < 2 || base > len(digits) {
		panicBase()
	}
	return base
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

// FormatUint32 works like FormatInt.
func FormatUint32(buf []byte, u uint32, base int) int {
	b := uint32(fixBase(base))
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
	return finish(buf, n, base > 0)
}

// FormatInt32 works like FormatInt.
func FormatInt32(buf []byte, i int32, base int) int {
	if i >= 0 {
		return FormatUint32(buf, uint32(i), base)
	}
	if len(buf) < 2 {
		panicBuffer()
	}
	n := FormatUint32(buf[1:], uint32(-i), base)
	if base > 0 {
		buf[0] = '-'
		n++
	} else {
		buf[0] = ' '
		buf[n] = '-'
	}
	return n
}

// FormatUint64 works like FormatInt.
func FormatUint64(buf []byte, u uint64, base int) int {
	b := uint64(fixBase(base))
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
	return finish(buf, n, base > 0)
}

// FormatInt64 works like FormatInt.
func FormatInt64(buf []byte, i int64, base int) int {
	if i >= 0 {
		return FormatUint64(buf, uint64(i), base)
	}
	if len(buf) < 2 {
		panicBuffer()
	}
	n := FormatUint64(buf[1:], uint64(-i), base)
	if base > 0 {
		buf[0] = '-'
		n++
	} else {
		buf[0] = ' '
		buf[n] = '-'
	}
	return n
}

// FormatUint works like FormatInt.
func FormatUint(buf []byte, u uint, base int) int {
	if unsafe.Sizeof(u) <= 4 {
		return FormatUint32(buf, uint32(u), base)
	}
	return FormatUint64(buf, uint64(u), base)
}

// FormatInt stores text representation of u in buf using 2 <= |base| <= 36.
// Unused portion of the buffer is filed with spaces.
// If base > 0 then formatted value is left-justified and FormatInt returns
// its length. If base < 0 then formatted value is right-justified and
// FormatInt returns offset to its first char.
func FormatInt(buf []byte, i, base int) int {
	if unsafe.Sizeof(i) <= 4 {
		return FormatInt32(buf, int32(i), base)
	}
	return FormatInt64(buf, int64(i), base)
}
