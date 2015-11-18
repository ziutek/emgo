package ppi

import (
	"arch/cortexm/exce"
	"mmio"
	"unsafe"

	"nrf51/internal"
)

// Event represents nRF51 event.
type Event struct {
	reg *mmio.U32
}

func GetEvent(te *internal.TasksEvents, n int) Event {
	return Event{&te.Events[n]}
}

// Happened returns happened flag (true if event happened).
func (e Event) Happened() bool {
	return e.reg.Load() != 0
}

// Clear clears happened flag for event.
func (e Event) Clear() {
	e.reg.Store(0)
}

func (e Event) temask() (*internal.TasksEvents, uint32) {
	ea := e.reg.Ptr()
	tea := ea & 0xfffff000
	n := (ea-tea)>>2 - 64
	return (*internal.TasksEvents)(unsafe.Pointer(tea)), 1 << n
}

// IntEnabled tells whether the event generates interrupts.
func (e Event) IntEnabled() bool {
	te, mask := e.temask()
	return te.IntEnSet.Load()&mask != 0
}

// EnableInt enables generating interrupts by event.
func (e Event) EnableInt() {
	te, mask := e.temask()
	te.IntEnSet.Store(mask)
}

// DisableInt disables generating interrupts by event.
func (e Event) DisableInt() {
	te, mask := e.temask()
	te.IntEnClr.Store(mask)
}

// IRQ returns exception number associated to event
func (e Event) IRQ() exce.Exce {
	te, _ := e.temask()
	return te.IRQ()
}

// Enabled tells whether the event is enabled, that is, it can update happened
// flag and will be routed to PPI. Only some peripherals (eg. RTC) implements
// this method.
func (e Event) Enabled() bool {
	te, mask := e.temask()
	return te.EvtEnSet.Load()&mask != 0
}

// Enable enables event as described in Enabled description. Only some
// peripherals (eg. RTC) implements this method. Note that EnableInt also
// enables update of happened flag but doesn't connect event to PPI.
func (e Event) Enable() {
	te, mask := e.temask()
	te.EvtEnSet.Store(mask)
}

// Disable disables event. Only some peripherals (eg. RTC) implements this
// method. Note that Disable doesn't affect interrupt generation if it is
// enabled using EnableInt.
func (e Event) Disable() {
	te, mask := e.temask()
	te.EvtEnClr.Store(mask)
}
