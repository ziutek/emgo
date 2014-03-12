package runtime

import "unsafe"

// Memset fills  the  first  n  bytes of the memory area
// pointed to by s with the constant byte c
func Memset(s unsafe.Pointer, c byte, n uint)