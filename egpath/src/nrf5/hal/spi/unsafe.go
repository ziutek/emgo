package spi

import (
	"unsafe"
)

func (d *Driver) AsyncWriteRead(out, in []byte) {
	d.AsyncWriteStringRead(*(*string)(unsafe.Pointer(&out)), in)
}

func (d *Driver) WriteRead(out, in []byte) int {
	return d.WriteStringRead(*(*string)(unsafe.Pointer(&out)), in)
}
