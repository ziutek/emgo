package bytes

import (
	"internal"
	"unsafe"
)

func Compare(s1, s2 []byte) int {
	n := len(s1)
	if n > len(s2) {
		n = len(s2)
	}
	p1 := unsafe.Pointer((*internal.SliceHeader)(unsafe.Pointer(&s1)).Data)
	p2 := unsafe.Pointer((*internal.SliceHeader)(unsafe.Pointer(&s2)).Data)
	ret := internal.Memcmp(p1, p2, uintptr(n))
	if ret != 0 || len(s1) == len(s2) {
		switch {
		case ret < 0:
			ret = -1
		case ret > 0:
			ret = 1
		}
		return ret
	}
	if len(s1) < len(s2) {
		return -1
	}
	return 1
}

func Equal(s1, s2 []byte) bool {
	if len(s1) != len(s2) {
		return false
	}
	p1 := unsafe.Pointer((*internal.SliceHeader)(unsafe.Pointer(&s1)).Data)
	p2 := unsafe.Pointer((*internal.SliceHeader)(unsafe.Pointer(&s2)).Data)
	return internal.Memcmp(p1, p2, uintptr(len(s1))) == 0
}

// Fill fills s with b.
func Fill(s []byte, b byte) {
	if len(s) == 0 {
		return
	}
	internal.Memset(unsafe.Pointer(&s[0]), b, uintptr(len(s)))
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
