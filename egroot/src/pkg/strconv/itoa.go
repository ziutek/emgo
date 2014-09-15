package strconv

import "unsafe"

const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

func panicIfZero(n int) {
	if n == 0 {
		panic("strconv: buffer too short")
	}
}

// Utoa32 converts u to string and returns offset to most significant digit.
func Utoa32(buf []byte, u uint32, base int) int {
	if base < 2 || base > len(digits) {
		panic("strconv: illegal base")
	}
	b := uint32(base)
	n := len(buf)
	for u != 0 {
		panicIfZero(n)
		n--
		newU := u / b
		buf[n] = digits[u-newU*b]
		u = newU
	}
	if n == len(buf) {
		panicIfZero(n)
		n--
		buf[n] = '0'
	}
	for k := 0; k < n; k++ {
		buf[k] = ' '
	}
	return n
}

// Itoa32 converts i to string and returns offset to most significant digit or
// sign.
func Itoa32(buf []byte, i int32, base int) int {
	if i >= 0 {
		return Utoa32(buf, uint32(i), base)
	}
	if len(buf) == 0 {
		panic("strconv: buffer too short")
	}
	n := Utoa32(buf[1:], uint32(-i), base)
	buf[n] = '-'
	return n
}

// Utoa64 converts u to string and returns offset to most significant digit.
func Utoa64(buf []byte, u uint64, base int) int {
	if base < 2 || base > len(digits) {
		panic("strconv: illegal base")
	}
	b := uint64(base)
	n := len(buf)
	for u != 0 {
		panicIfZero(n)
		n--
		newU := u / b
		buf[n] = digits[u-newU*b]
		u = newU
	}
	if n == len(buf) {
		panicIfZero(n)
		n--
		buf[n] = '0'
	}
	for k := 0; k < n; k++ {
		buf[k] = ' '
	}
	return n
}

// Itoa64 converts i to string and returns offset to most significant digit or
// sign.
func Itoa64(buf []byte, i int64, base int) int {
	if i >= 0 {
		return Utoa64(buf, uint64(i), base)
	}
	if len(buf) == 0 {
		panic("strconv: buffer too short")
	}
	n := Utoa64(buf[1:], uint64(-i), base)
	buf[n] = '-'
	return n
}

func Utoa(buf []byte, u uint, base int) int {
	if unsafe.Sizeof(u) <= 4 {
		return Utoa32(buf, uint32(u), base)
	}
	return Utoa64(buf, uint64(u), base)
}

func Itoa(buf []byte, i, base int) int {
	if unsafe.Sizeof(i) <= 4 {
		return Itoa32(buf, int32(i), base)
	}
	return Itoa64(buf, int64(i), base)
}
