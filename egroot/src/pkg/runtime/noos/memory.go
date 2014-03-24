package noos

import "unsafe"

var heap []byte

func allocBottom(sptr unsafe.Pointer, b []byte, n int, size, align uintptr) []byte

func allocTop(sptr unsafe.Pointer, b []byte, n int, size, align uintptr) []byte

func heapStack() []byte

func heapStackEnd() uintptr

func stackSize() uintptr

func panicMemory() {
	panic("not enough memory")
}