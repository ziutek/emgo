// Package stack allows allocate memory on stack using  C alloca function.
// Returned memory is zeroed.
package stack

import "unsafe"

func alloc(n int, size uintptr) unsafe.Pointer

func ints8(n int) []int8
func bytes(n int) []byte
func ints16(n int) []int16
func uints16(n int) []uint16
func ints32(n int) []int32
func uints32(n int) []uint32
func ints64(n int) []int64
func uints64(n int) []uint64
func ints(n int) []int
func uints(n int) []uint

func floats32(n int) []float32
func floats64(n int) []float64

func complexs64(n int) []complex64
func complexs128(n int) []complex128

func bools(n int) []bool

func interfaces(n int) []interface{}

func Alloc(n int, size uintptr) unsafe.Pointer {
	return alloc(n, size)
}

func Ints8(n int) []int8 {
	return ints8(n)
}

func Bytes(n int) []byte {
	return bytes(n)
}

func Ints16(n int) []int16 {
	return ints16(n)
}

func Uints16(n int) []uint16 {
	return uints16(n)
}

func Ints32(n int) []int32 {
	return ints32(n)
}

func Uints32(n int) []uint32 {
	return uints32(n)
}
