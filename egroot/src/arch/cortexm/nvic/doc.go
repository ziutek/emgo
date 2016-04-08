// Package nvic allows to control Nested Vectored Interrupt Controller. It does
// not expose NVIC registers directly but defines IRQ type with set of methods
// that allow manage external interrupts.
//
// This package can not manage system exceptions. Use scb package instead.
//
// NVIC combines level and pulse sensing of interrupt signals. It is important
// to understand this behavior to avoid subtle bugs in device drivers code.
//
// An interrupt can be in four states:
//	1. Inactive: the interrupt is not active and not pending.
//	2. Pending: the interrupt is waiting to be serviced by the processor.
//	3. Active: the interrupt is being serviced by the processor.
//	4. Active and pending: the interrupt is being serviced by the processor and
//	   there is a pending interrupt from the same source.
//
// In simple terms, it can be assumed that level sensing is used in inactive
// state, pulse sensing is used in active state. This behavior has two
// important consequences:
//
// If the interrupt signal remains active after return from the interrupt
// handler, the interrupt state becomes pending and the handler will be
// executed again.
//
// If the interrupt signal was deasserted and next asserted again in active
// state, the interrupt state changes to active and pending, the handler will
// be executed again after return.
package nvic
