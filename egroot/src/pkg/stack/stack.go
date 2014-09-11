// Package stack allows allocate memory on stack using C alloca function.
// Allocated memory is zeroed.
package stack

import "unsafe"

func alloc(n int, size uintptr) unsafe.Pointer

func ints(n int) []int
func uints(n int) []uint
func uintptrs(n int) []uintptr
func bools(n int) []bool
func interfaces(n int) []interface{}

func ints8(n int) []int8
func uints8(n int) []byte
func ints16(n int) []int16
func uints16(n int) []uint16
func ints32(n int) []int32
func uints32(n int) []uint32
func ints64(n int) []int64
func uints64(n int) []uint64

func floats32(n int) []float32
func floats64(n int) []float64

func complexs64(n int) []complex64
func complexs128(n int) []complex128

func Alloc(n int, size uintptr) unsafe.Pointer {
	return alloc(n, size)
}

func Ints(n int) []int {
	return ints(n)
}

func Uints(n int) []uint {
	return uints(n)
}

func Uintptrs(n int) []uintptr {
	return uintptrs(n)
}

func Interfaces(n int) []interface{} {
	return interfaces(n)
}

func Ints8(n int) []int8 {
	return ints8(n)
}

func Ints16(n int) []int16 {
	return ints16(n)
}

func Ints32(n int) []int32 {
	return ints32(n)
}

func Ints64(n int) []int64 {
	return ints64(n)
}

func Bytes(n int) []byte {
	return uints8(n)
}

func Uints16(n int) []uint16 {
	return uints16(n)
}

func Uints32(n int) []uint32 {
	return uints32(n)
}

func Uints64(n int) []uint64 {
	return uints64(n)
}

func Float32(n int) []float32 {
	return floats32(n)
}

func Float64(n int) []float64 {
	return floats64(n)
}

func Complexs64(n int) []complex64 {
	return complexs64(n)
}
