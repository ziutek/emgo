package noos

import (
	"unsafe"
)

var Heap = heap()

func heap() []byte

func allocBottom(sptr unsafe.Pointer, b []byte, n int, elSize, elAlign, sliAlign uintptr) []byte

func allocTop(sptr unsafe.Pointer, b []byte, n int, elSize, elAlign, sliAlign uintptr) []byte

func panicMemory() {
	panic("not enough memory for runtime initialisation")
}
