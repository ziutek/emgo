package sdcard

import (
	"reflect"
	"unsafe"
)

// Data should be used for data transfers. It ensures 8-byte alignment.
type Data []uint64

// Bytes returns d as []byte.
func (d Data) Bytes() []byte {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	h.Len *= 8
	h.Cap *= 8
	return *(*[]uint8)(unsafe.Pointer(&d))
}
