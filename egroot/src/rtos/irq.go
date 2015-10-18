package rtos

import "syscall"

// An IRQ represents exception/interrupt type.
type IRQ int

func (irq IRQ) Enable() error {
	return mkerror(syscall.SetIRQEna(int(irq), true))
}

func (irq IRQ) Disable() error {
	return mkerror(syscall.SetIRQEna(int(irq), false))
}

// IPrio represents IRQ priority.
type IPrio int

const (
	IPrioLowest  = IPrio(syscall.IRQPrioLowest)
	IPrioHighest = IPrio(syscall.IRQPrioHighest)
)

// Lower resturns true if priority p is lower than o.
func (p IPrio) Lower(o IPrio) bool {
	if IPrioLowest < IPrioHighest {
		return p < o
	}
	return p > o
}

// Higher resturns true if priority p is higher than o.
func (p IPrio) Higher(o IPrio) bool {
	if IPrioLowest > IPrioHighest {
		return p > o
	}
	return p < o
}

// SetPrio sets priority for irq.
func (irq IRQ) SetPrio(p IPrio) error {
	return mkerror(syscall.SetIRQPrio(int(irq), int(p)))
}

// UseHandler sets h as handler for irq.
func (irq IRQ) UseHandler(h func()) error {
	return mkerror(syscall.SetIRQHandler(int(irq), h))
}

// Status returns basic information about irq.
func (irq IRQ) Status() (prio IPrio, enabled bool, err error) {
	s, e := syscall.IRQStatus(int(irq))
	enabled = s < 0
	if enabled {
		s = -s - 1
	}
	prio = IPrio(s)
	err = mkerror(e)
	return
}
