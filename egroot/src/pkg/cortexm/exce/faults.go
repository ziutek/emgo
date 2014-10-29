package exce

import "unsafe"

// FaultStatusRegs represents Cortex-M3 memory-mapped registers that contains details
// about faults..
type FaultStatusRegs struct {
	MMS uint8   `C:"volatile"` // MemManage fault status
	BFS uint8   `C:"volatile"` // BusFalut status
	UFS uint16  `C:"volatile"` // UsageFault status
	HFS uint32  `C:"volatile"` // HardFault status
	_   uint32  `C:"volatile"`
	MMA uintptr `C:"volatile"` // MemManage fault address
	BFA uintptr `C:"volatile"` // BusFault address
	AFS uint32  `C:"volatile"` // Auxiliary fault Status
}

// FSR points to area of memory that contains fault status registers..
var FSR = (*FaultStatusRegs)(unsafe.Pointer(uintptr(0xe000ed28)))
