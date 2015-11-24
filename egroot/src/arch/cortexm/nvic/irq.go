// Package nvic allows to control Nested Vectored Interrupt Controller. It does
// not expose NVIC registers directly but defines IRQ type with set of methods
// that allow manage external interrupts.
//
// This package can not manage system exceptions. Use scb package instead.
package nvic

// IRQ represents Cortex-M external interrupt.
type IRQ byte

// Enabled returns true if handler for irq is enabled.
func (irq IRQ) Enabled() bool {
	return r.ISER.Bit(irq)
}

// Enable enables handler for e.
func (irq IRQ) Enable() {
	r.ISER.SetBit(irq)
}

// Disable disables handler for irq. To disable some class of exceptions in
// atomic way see PRIMASK, FAULTMASK, BASEPRI registers in cortexm package..
func (irq IRQ) Disable() {
	r.ICER.SetBit(irq)
}

// Pending returns true if irq is pending.
func (irq IRQ) Pending() bool {
	return r.ISPR.Bit(irq)
}

// SetPending generates irq.
func (irq IRQ) SetPending() {
	r.ISPR.SetBit(irq)
}

// ClearPending clears pending flag for irq.
func (irq IRQ) ClearPending() {
	r.ICPR.SetBit(irq)
}

// Active returns true if CPU is handling irq.
func (irq IRQ) Active() bool {
	return r.IABR.Bit(irq)
}

// Prio returns priority level for irq.
func (irq IRQ) Prio() int {
	return int(r.IPR.Byte(irq))
}

// SetPrio sets priority level for irq.
func (irq IRQ) SetPrio(p int) {
	r.IPR.SetByte(irq, byte(p))
}

// Trig generates irq. This function can be call from user level if
// USERSETMPEND bit is set in scb.CCR.
func (irq IRQ) Trig() {
	sti.Store(uint32(irq))
}
