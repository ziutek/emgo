package bytes

import (
	"builtin"
	"unsafe"
)

// Fill fills s with b.
func Fill(s []byte, b byte) {
	if len(s) == 0 {
		return
	}
	builtin.Memset(unsafe.Pointer(&s[0]), b, uintptr(len(s)))
}

// IndexByte returns the index of first c in s or -1 if there is no c in s.
func IndexByte(s []byte, c byte) int {
	for i, b := range s {
		if b == c {
			return i
		}
	}
	return -1
}

/*
// Index returns the index of first sep in s or -1 if there is no sep in s.
func Index(s, sep []byte) int {
	if len(sep) == 0 {
		return 0
	}
	if len(sep) == 1 {
		return IndexByte(sep[0])
	}
	...
}
*/
