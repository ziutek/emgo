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
