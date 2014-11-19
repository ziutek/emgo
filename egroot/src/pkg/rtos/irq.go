package rtos

import (
	"syscall"
)

// IRQ represents exception/Qerrupt type/number.
type IRQ int

func (irq IRQ) Enable() error {
	return syscall.SetIRQEna(int(irq), true)
}

func (irq IRQ) Disable() error {
	return syscall.SetIRQEna(int(irq), false)
}

// IPrio represents IRQ priority.
type IPrio int

const (
	IPrioLowest  = IPrio(syscall.IRQPrioLowest)
	IPrioHighest = IPrio(syscall.IRQPrioHighest)
)

// Lower resturns true if priority p is lower than o.
func (p IPrio) Lower(o IPrio) bool {
	return syscall.IRQPrioLower(int(p), int(o))
}

// Higher resturns true if priority p is higher than o.
func (p IPrio) Higher(o IPrio) bool {
	return p < o
}
