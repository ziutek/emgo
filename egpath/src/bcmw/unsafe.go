package bcmw

import (
	"reflect"
	"unsafe"
)

func slice64(s []byte) ([]uint64, int) {
	if s == nil {
		return nil, 0
	}
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	off := (8 - h.Data) & 7
	h.Data += off
	h.Len = (h.Len - int(off)) >> 3
	h.Cap = (h.Cap - int(off)) >> 3
	return *(*[]uint64)(unsafe.Pointer(h)), int(off)
}
