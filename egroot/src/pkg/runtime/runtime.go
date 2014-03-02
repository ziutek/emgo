package runtime

import "unsafe"

// Copy copies n bytes from the location pointed by src
// to the location pointed by dst. Locations can overlap.
func Copy(dst, src unsafe.Pointer, n uint)

// Memset fills  the  first  n  bytes of the memory area
// pointed to by s with the constant byte c
func Memset(s unsafe.Pointer, c byte, n uint)

// Panic is called by builtin panic function. If Panic returns panic executes
// infinite loop.
var Panic func(s string) = defaultPanic

func defaultPanic(s string) {
	// do nothing
}
