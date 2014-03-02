package irq

import (
	"sync"
	"mmio"
	"unsafe"
)

// Vector represents element of interrupt table. In case of Cortex-M this is
// simply func().
type Vector func()

// SysTable represents table of Cortex-M system exceptions. Reset, NMI and
// HardFault are always enabled so they should always be set to correct handler
// function.
//
// If you modified table that is currently used by CPU, you should call
// sync.Memory() to be sure that modifications take effect. Alternatively you
// can setup new table and make it active using SetActiveTable.
type SysTable struct {
	_ Vector `C:"__attribute__((aligned(32*4)))"`

	Reset     Vector
	NMI       Vector
	HardFault Vector

	MemFault   Vector
	BusFault   Vector
	UsageFault Vector
	_          Vector
	_          Vector
	_          Vector
	_          Vector
	SVCall     Vector
	DebugMon   Vector
	_          Vector
	PendSV     Vector
	SysTick    Vector
}

// Set sets handler for specified system IRQ.
func (t *SysTable) Set(irq IRQ, v Vector) {
	(*[16]Vector)(unsafe.Pointer(t))[irq] = v
}

// Slice() returns SysTable as slice of vectors.
func (t *SysTable) Slice() []Vector {
	return (*[16]Vector)(unsafe.Pointer(t))[:]
}

var vto = (*mmio.Reg32)(unsafe.Pointer(uintptr(0xe000ed08)))

// SetActiveTable instruct CPU to use t as vector table. t should be properly
// aligned. Minimum alignment is 32 words which is enough for up to 16 external
// interrupts. For more interrupts, adjust the alignment by rounding up to the
// next power of two. SyncActiveTable calls sync.Memory() before using t. 
func SetActiveTable(t []Vector) {
	sync.Memory()
	vto.Write(uint32(uintptr(unsafe.Pointer(&t[0]))))
}
