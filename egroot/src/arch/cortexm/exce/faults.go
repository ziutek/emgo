package exce

import "unsafe"

// FaultStatusRegs represents Cortex-M3 memory-mapped registers that contains
// details about faults.
type FaultStatusRegs struct {
	MMS uint8  // MemManage fault status
	BFS uint8  // BusFalut status
	UFS uint16 // UsageFault status
	HFS uint32 // HardFault status
	_   uint32
	MMA uintptr // MemManage fault address
	BFA uintptr // BusFault address
	AFS uint32  // Auxiliary fault Status
} //c:volatile

// FSR points to area of memory that contains fault status registers..
var FSR = (*FaultStatusRegs)(unsafe.Pointer(uintptr(0xe000ed28)))
