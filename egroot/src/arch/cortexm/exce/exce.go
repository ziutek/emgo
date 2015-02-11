package exce

import (
	"mmio"
	"unsafe"
)

// Exce represents Cortex-M exception. Its value is equal to Exception Number
// (not IRQ number) as defined by ARM in Cortex-M documentation. Lowest Exc
// value is 1 (Reset), highes value is 255 (externel interrupt #239).
type Exce byte

// Cortex-M system exceptions (with they default priority levels).
const (
	Reset      Exce = iota + 1 // prio -3 (fixed)
	NMI                        // prio -2 (fixed)
	HardFault                  // prio -1 (fixed)
	MemManage                  // prio 0
	BusFault                   // prio 1
	UsageFault                 // prio 2
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
const IRQ0 Exce = 16

var (
	ise  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e100)))
	ice  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e180)))
	isp  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e200)))
	icp  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e280)))
	iab  = (*bitReg)(unsafe.Pointer(uintptr(0xe000e300)))
	ip   = (*byteReg)(unsafe.Pointer(uintptr(0xe000e400)))
	ics  = mmio.PtrReg32(0xe000ed04)
	shp  = (*byteReg)(unsafe.Pointer(uintptr(0xe000ed18)))
	shcs = mmio.PtrReg32(0xe000ed24)
	sti  = mmio.PtrReg32(0xe000ef00)
)

// Enable enables handler for irq.
func (e Exce) Enable() {
	switch {
	case e >= MemManage && e <= UsageFault:
		shcs.SetBit(int(18 - UsageFault + e))
	case e >= IRQ0:
		ise.setBit(e - IRQ0)
	}
}

// Enabled returns true if handler for irq is enabled.
func (e Exce) Enabled() bool {
	switch {
	case e >= MemManage && e <= UsageFault:
		return shcs.Bit(int(18 - UsageFault + e))
	case e >= IRQ0:
		return ise.bit(e - IRQ0)
	}
	return true
}

// Disable disables handler for irq. Only handlers for MemManage, BusFault,
// UsageFault system exceptions and handlers for external interrupts can be
// disabled. In case of system exceptions, disabled handler means that
// HardFault will be used. There are SetPrimask, SetFaultmask, SetBasePrimask
// functions that can be used to disable handling of some class of exceptions
// in one step.
func (e Exce) Disable() {
	switch {
	case e >= MemManage && e <= UsageFault:
		shcs.ClearBit(int(18 - UsageFault + e))
	case e >= IRQ0:
		ice.setBit(e - IRQ0)
	}
}

// Pending returns true if irq is pending.
func (e Exce) Pending() bool {
	if e >= IRQ0 {
		return isp.bit(e - IRQ0)
	}
	switch e {
	case NMI:
		return ics.Bit(31)

	case MemManage:
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
func (e Exce) SetPending() {
	if e >= IRQ0 {
		isp.setBit(e - IRQ0)
		return
	}
	switch e {
	case NMI:
		ics.SetBit(31)

	case MemManage:
		shcs.SetBit(13)

	case BusFault:
		shcs.SetBit(14)

	case UsageFault:
		shcs.SetBit(12)

	case SVCall:
		shcs.SetBit(15)

	case PendSV:
		ics.Store(1 << 28)

	case SysTick:
		ics.Store(1 << 26)
	}
}

// ClearPending cancel pending irq.
func (e Exce) ClearPending() {
	if e >= IRQ0 {
		icp.setBit(e - IRQ0)
		return
	}
	switch e {
	case MemManage:
		shcs.ClearBit(13)

	case BusFault:
		shcs.ClearBit(14)

	case UsageFault:
		shcs.ClearBit(12)

	case SVCall:
		shcs.ClearBit(15)

	case PendSV:
		ics.Store(1 << 27)

	case SysTick:
		ics.Store(1 << 25)
	}
}

// Active returns true if CPU is handling irq.
func (e Exce) Active() bool {
	if e >= IRQ0 {
		return iab.bit(e - IRQ0)
	}
	switch e {
	case MemManage:
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
func (e Exce) Prio() Prio {
	switch {
	case e >= MemManage && e < IRQ0:
		return Prio(shp.byte(e - MemManage))
	case e >= IRQ0:
		return Prio(ip.byte(e - IRQ0))
	}
	return Highest
}

func (e Exce) SetPrio(p Prio) {
	switch {
	case e >= MemManage && e < IRQ0:
		shp.setByte(e-MemManage, byte(p))

	case e >= IRQ0:
		ip.setByte(e-IRQ0, byte(p))

	default:
		panic("can't set priority for exception < MemFault")
	}
}

// Trig generates irq. Only external interrupts can be generated this way.
// This function can be call from user level if USERSETMPEND bit in CCR is set.
func (e Exce) Trig() {
	if e < IRQ0 {
		return
	}
	sti.Store(uint32(e - IRQ0))
}

// Pending return true if any exception other than NMI or fault is pending
func Pending() bool {
	return ics.Bit(22)
}

// VecPending returns the number of the highest priority pending enabled
// exception. 0 means no pending exceptions. Returned value includes the
// effect of the BASEPRI and FAULTMASK (but not PRIMASK) registers.
func VecPending() Exce {
	return Exce(ics.Load() >> 12)
}

// VecActive returns the number of active exception or 0 for thread mode.
func VecActive() Exce {
	return Exce(ics.Load())
}
