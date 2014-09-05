package strconv

import (
	"io"
	"log"
	"unsafe"
)

const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

func panicIfZero(n int) {
	if n == 0 {
		log.Panic("strconv: buffer too short")
	}
}

// Utoa32 converts u to string and returns offset to most significant digit.
func Utoa32(buf []byte, u uint32, base int) int {
	if base < 2 || base > len(digits) {
		log.Panic("strconv: illegal base")
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
		log.Panic("strconv: buffer too short")
	}
	n := Utoa32(buf[1:], uint32(-i), base)
	buf[n] = '-'
	return n
}

func WriteInt32(w io.Writer, i int32, base int) (int, error) {
	var buf [11]byte
	first := Itoa32(buf[:], i, base)
	return w.Write(buf[first:])
}

// Utoa64 converts u to string and returns offset to most significant digit.
func Utoa64(buf []byte, u uint64, base int) int {
	if base < 2 || base > len(digits) {
		log.Panic("strconv: illegal base")
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

func WriteUint64(w io.Writer, u uint64, base int) (int, error) {
	var buf [20]byte
	first := Utoa64(buf[:], u, base)
	return w.Write(buf[first:])
}

// Itoa64 converts i to string and returns offset to most significant digit or
// sign.
func Itoa64(buf []byte, i int64, base int) int {
	if i >= 0 {
		return Utoa64(buf, uint64(i), base)
	}
	if len(buf) == 0 {
		log.Panic("strconv: buffer too short")
	}
	n := Utoa64(buf[1:], uint64(-i), base)
	buf[n] = '-'
	return n
}

func WriteInt64(w io.Writer, i int64, base int) (int, error) {
	var buf [21]byte
	first := Itoa64(buf[:], i, base)
	return w.Write(buf[first:])
}

func Itoa(buf []byte, i, base int) int {
	if unsafe.Sizeof(i) <= 4 {
		return Itoa32(buf, int32(i), base)
	}
	return Itoa64(buf, int64(i), base)

}
