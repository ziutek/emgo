package noos

import "unsafe"

var heap []byte

func alloc(sptr unsafe.Pointer, b []byte, n int, size uintptr) []byte

func sliceU8(p unsafe.Pointer, n uint) []byte

func sliceU16(p unsafe.Pointer, n uint) []byte

func sliceU32(p unsafe.Pointer, n uint) []byte

func heapStack() []byte

func heapSize() uintptr

func panicMemory() {
	panic("not enough memory")
}