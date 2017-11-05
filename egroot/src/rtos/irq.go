package rtos

import "syscall"

// An IRQ represents an interrupt request number.
// Cortex-M: IRQ represents external interrupt request.
type IRQ int

// Enable enables handling of the IRQ.
func (irq IRQ) Enable() error {
	return mkerror(syscall.SetIRQEna(int(irq), true))
}

// Disable disables handling of the IRQ.
func (irq IRQ) Disable() error {
	return mkerror(syscall.SetIRQEna(int(irq), false))
}

// IRQPrio represents IRQ priority.
type IRQPrio int

const (
	// IRQPrioLowest is the lowest IRQ priority.
	IRQPrioLowest = IRQPrio(syscall.IRQPrioLowest)

	// IRQPrioHighest is the highest IRQ priority.
	IRQPrioHighest = IRQPrio(syscall.IRQPrioHighest)

	// IRQPrioNum is the number of priority levels.
	IRQPrioNum = IRQPrio(syscall.IRQPrioNum)

	// IRQPrioStep if added to priority increases it to next, highest level. In
	// many cases number of effective levels is less than IRQPrioNum and adding
	// one step to priority does not guarantee highest effective level.
	IRQPrioStep = IRQPrio(syscall.IRQPrioStep)

	// SyscallPrio is the IRQ priority level equal to priority of IRQ/exception
	// used to implement system calls. Usually, to use system calls in interrupt
	// handler its IRQ must have priority lower than SyscallPrio.
	SyscallPrio = IRQPrio(syscall.SyscallPrio)
)

// Lower resturns true if priority p is lower than o.
func (p IRQPrio) Lower(o IRQPrio) bool {
	if IRQPrioLowest < IRQPrioHighest {
		return p < o
	}
	return p > o
}

// Higher resturns true if priority p is higher than o.
func (p IRQPrio) Higher(o IRQPrio) bool {
	if IRQPrioLowest > IRQPrioHighest {
		return p > o
	}
	return p < o
}

// SetPrio sets priority for the IRQ.
func (irq IRQ) SetPrio(p IRQPrio) error {
	return mkerror(syscall.SetIRQPrio(int(irq), int(p)))
}

// UseHandler sets h as handler for the IRQ. It can be not supported by some
// architectures or when vector table is located in ROM/Flash.
func (irq IRQ) UseHandler(h func()) error {
	return mkerror(syscall.SetIRQHandler(int(irq), h))
}

// Status returns basic information about the IRQ.
func (irq IRQ) Status() (prio IRQPrio, enabled bool, err error) {
	s, e := syscall.IRQStatus(int(irq))
	enabled = s < 0
	if enabled {
		s = -s - 1
	}
	prio = IRQPrio(s)
	err = mkerror(e)
	return
}

// Trigger allows to trigger the IRQ by software.
func (irq IRQ) Trigger() {
	syscall.TriggerIRQ(int(irq))
}
