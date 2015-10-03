package strconv

import "unsafe"

const (
	intSize = unsafe.Sizeof(int(0))
	ptrSize = unsafe.Sizeof(uintptr(0))
)