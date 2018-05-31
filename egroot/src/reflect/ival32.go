// +build !amd64

package reflect

type ival struct {
	ptr uintptr
	w   uint32
	dw  uint64
}
