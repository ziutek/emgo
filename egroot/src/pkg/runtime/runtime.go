package runtime

import (
	"unsafe"
)

// Copy copies n bytes from the location pointed by src
// to the location pointed by dst. Locations can overlap.
func Copy(dst, src unsafe.Pointer, n uint)

func Panic(s string) {
	for {
	}
}