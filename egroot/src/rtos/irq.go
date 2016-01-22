package rtos

import "syscall"

// An IRQ represents interrupt type.
// Cortex-M: IRQ represents external interrupt.
type IRQ int

func (irq IRQ) Enable() error {
	return mkerror(syscall.SetIRQEna(int(irq), true))
}

func (irq IRQ) Disable() error {
	return mkerror(syscall.SetIRQEna(int(irq), false))
}

// IRQPrio represents IRQ priority.
type IRQPrio int

const (
	IRQPrioLowest  = IRQPrio(syscall.IRQPrioLowest)
	IRQPrioHighest = IRQPrio(syscall.IRQPrioHighest)
	IRQPrioStep    = IRQPrio(syscall.IRQPrioStep)
	IRQPrioNum     = IRQPrio(syscall.IRQPrioNum)
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

// SetPrio sets priority for irq.
func (irq IRQ) SetPrio(p IRQPrio) error {
	return mkerror(syscall.SetIRQPrio(int(irq), int(p)))
}

// UseHandler sets h as handler for irq. It can be not supported by some
// architectures or when vector table is located in ROM/Flash. 
func (irq IRQ) UseHandler(h func()) error {
	return mkerror(syscall.SetIRQHandler(int(irq), h))
}

// Status returns basic information about irq.
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
