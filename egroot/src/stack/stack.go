// Package stack allows allocate memory on stack using C alloca function.
// Allocated memory is zeroed.
package stack

import "unsafe"

//c:inline
func Alloc(n int, size uintptr) uintptr

//c:inline
func Ints(n int) []int

//c:inline
func Uints(n int) []uint

//c:inline
func Uintptrs(n int) []uintptr

//c:inline
func Pointers(n int) []unsafe.Pointer

//c:inline
func Interfaces(n int) []interface{}

//c:inline
func Ints8(n int) []int8

//c:inline
func Ints16(n int) []int16

//c:inline
func Ints32(n int) []int32

//c:inline
func Ints64(n int) []int64

//c:inline
func Bytes(n int) []byte

//c:inline
func Uints16(n int) []uint16

//c:inline
func Uints32(n int) []uint32

//c:inline
func Uints64(n int) []uint64

//c:inline
func Floats32(n int) []float32

//c:inline
func Floats64(n int) []float64

//c:inline
func Complexs64(n int) []complex64

//c:inline
func Complexs128(n int) []complex64

//c:inline
func Bools(n int) []bool

