package nvic

import "arch/cortexm/scb"

// Exce represents Cortex-M exception. Its value is equal to IRQ number as
// defined by ARM in Cortex-M documentation. Lowest Exce value is -15 (Reset),
// highes value is 239.
type Exce int

// Cortex-M system exceptions.
const (
	Reset      Exce = -15
	NMI        Exce = -14
	HardFault  Exce = -13
	MemManage  Exce = -12
	BusFault   Exce = -11
	UsageFault Exce = -10
	SVC        Exce = -5
	DebugMon   Exce = -4
	PendSV     Exce = -2
	SysTick    Exce = -1
)

// First external interrupt
const IRQ0 Exce = 0

// Enabled returns true if handler for e is enabled.
func (e Exce) Enabled() bool {
	if e >= IRQ0 {
		return r.ISER.Bit(e)
	}
	switch e {
	case MemManage:
		return scb.SHCSR.Bits(scb.MEMFAULTENA) != 0
	case BusFault:
		return scb.SHCSR.Bits(scb.BUSFAULTENA) != 0
	case UsageFault:
		return scb.SHCSR.Bits(scb.USGFAULTENA) != 0
	}
	return true
}

// Enable enables handler for e.
func (e Exce) Enable() {
	if e >= IRQ0 {
		r.ISER.SetBit(e)
		return
	}
	switch e {
	case MemManage:
		scb.SHCSR.SetBits(scb.MEMFAULTENA)
	case BusFault:
		scb.SHCSR.SetBits(scb.BUSFAULTENA)
	case UsageFault:
		scb.SHCSR.SetBits(scb.USGFAULTENA)
	}
}

// Disable disables handler for e. In case of system exceptions, disable
// handler means that HardFault handler will be used instead. To disable some
// class of exceptions in atomic way use PRIMASK, FAULTMASK, BASEPRI registers
// (see functions in arch/cortexm package).
func (e Exce) Disable() {
	if e >= IRQ0 {
		r.ICER.SetBit(e)
	}
	switch e {
	case MemManage:
		scb.SHCSR.ClearBits(scb.MEMFAULTENA)
	case BusFault:
		scb.SHCSR.ClearBits(scb.BUSFAULTENA)
	case UsageFault:
		scb.SHCSR.ClearBits(scb.USGFAULTENA)
	}
}

// Pending returns true if e is pending.
func (e Exce) Pending() bool {
	if e >= IRQ0 {
		return r.ISPR.Bit(e)
	}
	switch e {
	case NMI:
		return scb.ICSR.Bits(scb.NMIPENDSET) != 0
	case MemManage:
		return scb.SHCSR.Bits(scb.MEMFAULTPENDED) != 0
	case BusFault:
		return scb.SHCSR.Bits(scb.BUSFAULTPENDED) != 0
	case UsageFault:
		return scb.SHCSR.Bits(scb.USGFAULTPENDED) != 0
	case SVC:
		return scb.SHCSR.Bits(scb.SVCALLPENDED) != 0
	case PendSV:
		return scb.ICSR.Bits(scb.PENDSVSET) != 0
	case SysTick:
		return scb.ICSR.Bits(scb.PENDSTSET) != 0
	}
	return false
}
