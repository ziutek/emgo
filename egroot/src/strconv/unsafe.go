package strconv

import "unsafe"

const (
	intSize = unsafe.Sizeof(int(0))
	ptrSize = unsafe.Sizeof(uintptr(0))
)

func ParseUint32(s []byte, base int) (uint32, error) {
	return ParseStringUint32(*(*string)(unsafe.Pointer(&s)), base)
}

func ParseUint64(s []byte, base int) (uint64, error) {
	return ParseStringUint64(*(*string)(unsafe.Pointer(&s)), base)
}

func ParseUint(s []byte, base int) (uint, error) {
	return ParseStringUint(*(*string)(unsafe.Pointer(&s)), base)
}
