package irq

import (
	"mmio"
	"unsafe"
)

// IRQ represents Cortex-M exception. Its value is equal to Exception Number
// (not IRQ number) as defined by ARM in Cortex-M documentation. Lowest IRQ
// value is 1 (Reset), highes value is 255 (externel interrupt #239).
type IRQ byte

// Cortex-M system exceptions (with they default priority levels).
const (
	Reset      IRQ = iota + 1 // prio -3 (fixed)
	NMI                       // prio -2 (fixed)
	HardFault                 // prio -1 (fixed)
	MemFault                  // prio 0
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
	ics  = mmio.NewReg32(0xe000ed04)
	shp  = (*byteReg)(unsafe.Pointer(uintptr(0xe000ed18)))
	shcs = mmio.NewReg32(0xe000ed24)
	sti  = mmio.NewReg32(0xe000ef00)
)

// Enable enables handler for irq.
func (irq IRQ) Enable() {
	switch {
	case irq >= MemFault && irq <= UsageFault:
		shcs.SetBit(int(18 - UsageFault + irq))
	case irq >= Ext0:
		ise.setBit(irq - Ext0)
	}
}

// Enabled returns true if handler for irq is enabled.
func (irq IRQ) Enabled() bool {
	switch {
	case irq >= MemFault && irq <= UsageFault:
		return shcs.Bit(int(18 - UsageFault + irq))
	case irq >= Ext0:
		return ise.bit(irq - Ext0)
	}
	return true
}

// Disable disables handler for irq. Only handlers for MemManage, BusFault,
// UsageFault system exceptions and handlers for external interrupts can be
// disabled. In case of system exceptions, disabled handler means that
// HardFault will be used. There are SetPrimask, SetFaultmask, SetBasePrimask
// functions that can be used to disable handling of some class of exceptions
// in one step.
func (irq IRQ) Disable() {
	switch {
	case irq >= MemFault && irq <= UsageFault:
		shcs.ClearBit(int(18 - UsageFault + irq))
	case irq >= Ext0:
		ice.setBit(irq - Ext0)
	}
}

// Pending returns true if irq is pending.
func (irq IRQ) Pending() bool {
	if irq >= Ext0 {
		return isp.bit(irq - Ext0)
	}
	switch irq {
	case NMI:
		return ics.Bit(31)

	case MemFault:
		return shcs.Bit(13)

	case BusFault:
		return shcs.Bit(14)

	case UsageFault:
		return shcs.Bit(12)

	case SVCall:
		return shcs.Bit(15)

	case PendSV:
		return ics.Bit(28)

	case SysTick:
		return ics.Bit(26)
	}
	return false
}

// SetPending generates irq.
func (irq IRQ) SetPending() {
	if irq >= Ext0 {
		isp.setBit(irq - Ext0)
		return
	}
	switch irq {
	case NMI:
		ics.SetBit(31)

	case MemFault:
		shcs.SetBit(13)

	case BusFault:
		shcs.SetBit(14)

	case UsageFault:
		shcs.SetBit(12)

	case SVCall:
		shcs.SetBit(15)

	case PendSV:
		ics.Write(1 << 28)

	case SysTick:
		ics.Write(1 << 26)
	}
}

// ClearPending cancel pending irq.
func (irq IRQ) ClearPending() {
	if irq >= Ext0 {
		icp.setBit(irq - Ext0)
		return
	}
	switch irq {
	case MemFault:
		shcs.ClearBit(13)

	case BusFault:
		shcs.ClearBit(14)

	case UsageFault:
		shcs.ClearBit(12)

	case SVCall:
		shcs.ClearBit(15)

	case PendSV:
		ics.Write(1 << 27)

	case SysTick:
		ics.Write(1 << 25)
	}
}

// Active returns true if CPU is handling irq.
func (irq IRQ) Active() bool {
	if irq >= Ext0 {
		return iab.bit(irq - Ext0)
	}
	switch irq {
	case MemFault:
		return shcs.Bit(0)

	case BusFault:
		return shcs.Bit(1)

	case UsageFault:
		return shcs.Bit(3)

	case SVCall:
		return shcs.Bit(7)

	case DebugMon:
		return shcs.Bit(8)

	case PendSV:
		return shcs.Bit(10)

	case SysTick:
		return shcs.Bit(11)
	}
	return false
}

// Priority returns priority level for irq. It returns Highest for Reset,
// NMI and HardFault but they real priority values are fixed at -3, -2, -1
// respectively.
func (irq IRQ) Prio() Prio {
	switch {
	case irq >= MemFault && irq < Ext0:
		return Prio(shp.byte(irq - MemFault))
	case irq >= Ext0:
		return Prio(ip.byte(irq - Ext0))
	}
	return Highest
}

func (irq IRQ) SetPrio(p Prio) {
	switch {
	case irq >= MemFault && irq < Ext0:
		shp.setByte(irq-MemFault, byte(p))

	case irq >= Ext0:
		ip.setByte(irq-Ext0, byte(p))

	default:
		panic("can't set priority for irq < MemFault")
	}
}

// Trig generates irq. Only external interrupts can be generated this way.
// This function can be call from user level if USERSETMPEND bit in CCR is set.
func (irq IRQ) Trig() {
	if irq < Ext0 {
		return
	}
	sti.Write(uint32(irq - Ext0))
}

// Pending return true if any interrupt other than NMI or fault is pending
func Pending() bool {
	return ics.Bit(22)
}

// VecPending returns the number of the highest priority pending enabled
// exception. 0 means no pending exceptions. Returned value includes the
// effect of the BASEPRI and FAULTMASK (but not PRIMASK) registers.
func VecPending() IRQ {
	return IRQ(ics.Read() >> 12)
}

// VecActive returns the number of active exception or 0 for thread mode.
func VecActive() IRQ {
	return IRQ(ics.Read())
}
