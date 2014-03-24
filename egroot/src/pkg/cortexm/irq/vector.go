package irq

import (
	"mmio"
	"sync/barrier"
	"unsafe"
)

// Vector represents element of interrupt table.
type Vector struct {
	handler func()
}

// VectorFor returns (implementation specific) Vector that correspods to the
// handler.
func VectorFor(handler func()) Vector {
	return Vector{handler}
}

var vto = mmio.NewReg32(0xe000ed08)

// UseTable instruct CPU to use vt as vector table. vt should be properly
// aligned. Minimum alignment is 32 words which is enough for up to 16 external
// interrupts. For more interrupts, adjust the alignment by rounding up to the
// next power of two. UseTable calls barrier.Memory() before seting vt as
// current table. UseTable is not available for Cortex-M0.
func UseTable(vt []Vector) {
	barrier.Memory()
	vto.Write(uint32(uintptr(unsafe.Pointer(&vt[0]))))
}
