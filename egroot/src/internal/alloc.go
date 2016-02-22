package internal

import "unsafe"

const (
	MemNoDMA     = 1 << iota
	MemStack     = -1
)

// Alloc is used for dynamic memory allocation. It alloates memory for n elements of
// specified size and alignment.
var Alloc func(n int, esize, ealign uintptr) unsafe.Pointer
