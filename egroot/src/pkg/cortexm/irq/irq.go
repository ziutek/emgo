package irq

import "unsafe"

// IRQ represents Cortex-M exception. Its value is equal to Exception
// Number (not IRQ number) as defined by ARM in Cortex-M documentation.
// Lowest IRQ value is 1 (Reset), highes value is 255 (externel
// interrupt #239).
type IRQ byte

// Cortex-M system exceptions (with they default priority levels).
const (
	Reset      IRQ = iota + 1 // prio -3 (fixed)
	NMI                       // prio -2 (fixed)
	HardFault                 // prio -1 (fixed)
	MemManage                 // prio 0
	BusFault                  // prio 1
	UsageFault                // prio 2
	_
	_
	_
	_
	SVCall   // prio 3
	DebugMon // prio 4
	_
	PendSV  // prio 5
	SysTick // prio 6
)

// First external interrupt
const Ext0 IRQ = 16

var (
	ise  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e100)))
	ice  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e180)))
	isp  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e200)))
	icp  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e280)))
	iab  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e300)))
	ip   = (*byteReg)(unsafe.Pointer(uintptr(0xe000e400)))
	shcs = (*wordReg)(unsafe.Pointer(uintptr(0xe000ed24)))
	sti  = (*wordReg)(unsafe.Pointer(uintptr(0xe000eF00)))
)

// Enable enables handler for irq.
func (irq IRQ) Enable() {
	switch {
	case irq >= MemManage && irq <= UsageFault:
		shcs.setBit(18 - UsageFault + irq)
	case irq >= Ext0:
		ise.setBit(irq - Ext0)
	}
}

// Enabled returns true if handler for irq is enabled.
func (irq IRQ) Enabled() bool {
	switch {
	case irq >= MemManage && irq <= UsageFault:
		return shcs.bit(18 - UsageFault + irq)
	case irq >= Ext0:
		return ise.bit(irq - Ext0)
	}
	return true
}

// Disable disables handler for irq. Only handlers for MemManage,
// BusFault, UsageFault system exceptions and handlers for external
// interrupts can be disabled. In case of system exceptions, disabled
// handler means that HardFault will be used. There are Primask,
// Faultmask, BasePrimask functions that can be used to disable
// handling of some class of exceptions in one step.
func (irq IRQ) Disable() {
	switch {
	case irq >= MemManage && irq <= UsageFault:
		shcs.clearBit(18 - UsageFault + irq)
	case irq >= Ext0:
		ice.setBit(irq - Ext0)
	}
}

// Pending returns true if irq is pending.
func (irq IRQ) Pending() bool {
	if irq >= Ext0 {
		return isp.bit(irq - Ext0)
	}
	panic("can't get pending flag for system exception")
}

// ClearPending cancel pending irq.
func (irq IRQ) ClearPending() {
	if irq >= Ext0 {
		icp.setBit(irq - Ext0)
	}
}

// Active returns true if CPU is handling irq.
func (irq IRQ) Active() bool {
	if irq >= Ext0 {
		return iab.bit(irq - Ext0)
	}
	panic("can't get active flag for system exception")
}

// Prio represents Cortex-M setable interrupt priority.
type Prio byte

const (
	Highest Prio = 0
	Lowest  Prio = 255
)

// SetPriority sets priority level for irq
func (irq IRQ) SetPriority(prio Prio) {
	switch {
	case irq >= MemManage && irq < Ext0:
		// TODO
	case irq >= Ext0:
		ip.r[irq-Ext0] = prio
	}
}

// Priority returns priority level for irq. It returns Highest
// for Reset, NMI and HardFault but they real priority values
// are fixed at -3, -2 and  -1 respectively.
func (irq IRQ) Priority() Prio {
	switch {
	case irq >= MemManage && irq < Ext0:
		// TODO
	case irq >= Ext0:
		return ip.r[irq-Ext0]
	}
	return Highest
}

// Trig generates irq. Only external interrupts can be
// generated this way.
func (irq IRQ) Trig() {
	if irq < Ext0 {
		return
	}
	sti.r = uint32(irq-Ext0)
}