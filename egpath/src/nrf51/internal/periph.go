package internal

import (
	"arch/cortexm/exce"
	"unsafe"
	"mmio"
)

// Base addresses for peripherals:
const (
	BaseAPB uintptr = 0x40000000 // accessed by APB,
	BaseAHB uintptr = 0x50000000 // accessed by AHB.
)

// TasksEvents should be the first field on any Periph struct.
// It takes 0x400 bytes of memory.
type TasksEvents struct {
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
func (te *TasksEvents) IRQ() exce.Exce {
	addr := uintptr(unsafe.Pointer(te))
	return exce.IRQ0 + exce.Exce((addr-BaseAPB)>>12)
}
