package bcmw

import (
	"reflect"
	"unsafe"
)

func uint64slice(s []byte) ([]byte, []uint64, []byte) {
	if s == nil {
		return nil, nil, nil
	}
	h := *(*reflect.SliceHeader)(unsafe.Pointer(&s))
	off := (8 - h.Data) & 7
	h.Data += off
	h.Len = (h.Len - int(off)) >> 3
	h.Cap = (h.Cap - int(off)) >> 3
	return s[:off], *(*[]uint64)(unsafe.Pointer(&h)), s[int(off)+h.Len*8:]
}
