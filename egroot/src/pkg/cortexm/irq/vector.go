package irq

import (
	"mmio"
	"sync/barrier"
	"unsafe"
)

// Vector represents element of interrupt table. In case of Cortex-M this is
// simply func().
type vector func()

// SysTable represents table of Cortex-M system exceptions. Reset, NMI and
// HardFault are always enabled so they should always be set to correct handler
// function.
//
// If you modified table that is currently used by CPU, you should call
// sync.Memory() to be sure that modifications take effect. Alternatively you
// can setup new table and make it active using SetActiveTable.
type SysTable struct {
	_ vector `C:"__attribute__((aligned(32*4)))"`

	Reset     vector
	NMI       vector
	HardFault vector

	MemFault   vector
	BusFault   vector
	UsageFault vector
	_          vector
	_          vector
	_          vector
	_          vector
	SVCall     vector
	DebugMon   vector
	_          vector
	PendSV     vector
	SysTick    vector
}

// Vector returns (implementation specific) value that correspods to the
// handler and can be an element of interrupt table.
func Vector(handler func()) vector {
	return handler
}

// Set sets handler for specified system IRQ.
func (t *SysTable) Set(irq IRQ, f func()) {
	(*[16]vector)(unsafe.Pointer(t))[irq] = f
}

// Slice() returns SysTable as slice of vectors.
func (t *SysTable) Slice() []vector {
	return (*[16]vector)(unsafe.Pointer(t))[:]
}

var vto = mmio.NewReg32(0xe000ed08)

// SetActiveTable instruct CPU to use t as vector table. t should be properly
// aligned. Minimum alignment is 32 words which is enough for up to 16 external
// interrupts. For more interrupts, adjust the alignment by rounding up to the
// next power of two. SetActiveTable calls sync.Memory() before seting t as 
// active table.
func SetActiveTable(t []vector) {
	barrier.Memory()
	vto.Write(uint32(uintptr(unsafe.Pointer(&t[0]))))
}
