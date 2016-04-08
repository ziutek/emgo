// Package stack allows allocate memory on stack using C alloca function.
// Allocated memory is zeroed.
package stack

import "unsafe"

//c:static inline
func Alloc(n int, size uintptr) uintptr

//c:static inline
func Ints(n int) []int

//c:static inline
func Uints(n int) []uint

//c:static inline
func Uintptrs(n int) []uintptr

//c:static inline
func Pointers(n int) []unsafe.Pointer

//c:static inline
func Interfaces(n int) []interface{}

//c:static inline
func Ints8(n int) []int8

//c:static inline
func Ints16(n int) []int16

//c:static inline
func Ints32(n int) []int32

//c:static inline
func Ints64(n int) []int64

//c:static inline
func Bytes(n int) []byte

//c:static inline
func Uints16(n int) []uint16

//c:static inline
func Uints32(n int) []uint32

//c:static inline
func Uints64(n int) []uint64

//c:static inline
func Float32(n int) []float32

//c:static inline
func Float64(n int) []float64

//c:static inline
func Complexs64(n int) []complex64

//c:static inline
func Complexs128(n int) []complex64
