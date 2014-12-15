// Package stack allows allocate memory on stack using C alloca function.
// Allocated memory is zeroed.
package stack

import "unsafe"

func Alloc(n int, size uintptr) unsafe.Pointer

func Ints(n int) []int
func Uints(n int) []uint
func Uintptrs(n int) []uintptr
func Pointers(n int) []unsafe.Pointer
func Interfaces(n int) []interface{}

func Ints8(n int) []int8
func Ints16(n int) []int16
func Ints32(n int) []int32
func Ints64(n int) []int64

func Bytes(n int) []byte
func Uints16(n int) []uint16
func Uints32(n int) []uint32
func Uints64(n int) []uint64

func Float32(n int) []float32
func Float64(n int) []float64

func Complexs64(n int) []complex64
func Complexs128(n int) []complex64
