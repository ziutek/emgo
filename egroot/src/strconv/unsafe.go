package strconv

import (
	"io"
	"unsafe"
)

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

func ParseInt32(s []byte, base int) (int32, error) {
	return ParseStringInt32(*(*string)(unsafe.Pointer(&s)), base)
}

func ParseInt64(s []byte, base int) (int64, error) {
	return ParseStringInt64(*(*string)(unsafe.Pointer(&s)), base)
}

func ParseInt(s []byte, base int) (int, error) {
	return ParseStringInt(*(*string)(unsafe.Pointer(&s)), base)
}

func WriteBytes(w io.Writer, s []byte, width int, pad rune) (int, error) {
	return WriteString(w, *(*string)(unsafe.Pointer(&s)), width, pad)
}
