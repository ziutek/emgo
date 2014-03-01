package irq

import "unsafe"

type Vector func()

// SysTable represents table of Cortex-M system exceptions. Reset, NMI
// and HardFault are always enabled so they should always be set to
// correct handler function.
//
// If you modified table that is currently used by CPU, you should call
// sync.Sync() to be sure that modifications take effect.
type SysTable struct {
	_          Vector `C:"__attribute__((aligned(16*4)))"`

	Reset      Vector
	NMI        Vector
	HardFault  Vector
	
	MemManage  Vector
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
