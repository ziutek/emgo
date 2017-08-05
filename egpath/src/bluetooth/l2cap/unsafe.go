package l2cap

import (
	"unsafe"
)

// Write works like WriteString.
func (f *BLEFAR) Write(s []byte) (n int, err error) {
	return f.WriteString(*(*string)(unsafe.Pointer(&s)))
}
