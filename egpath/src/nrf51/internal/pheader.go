package internal

import (
	"arch/cortexm/exce"
	"mmio"
	"unsafe"
)

// Pheader should be the first field on any Periph struct.
// It takes 0x400 bytes of memory.
type Pheader struct {
	Tasks    [32]mmio.U32
	_        [32]mmio.U32
	Events   [32]mmio.U32
	_        [32]mmio.U32
	Shorts   mmio.U32
	_        [64]mmio.U32
	IntEnSet mmio.U32
	IntEnClr mmio.U32
	_        [14]mmio.U32
	EvtEnSet mmio.U32
	EvtEnClr mmio.U32
	_        [45]mmio.U32
}

// IRQ returns exception number associated to events.
func (ph *Pheader) IRQ() exce.Exce {
	addr := uintptr(unsafe.Pointer(ph))
	return exce.IRQ0 + exce.Exce((addr-BaseAPB)>>12)
}
