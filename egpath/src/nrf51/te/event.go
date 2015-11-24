package te

import (
	"mmio"
	"unsafe"

	"arch/cortexm/nvic"

	"nrf51/internal"
)

// Event represents nRF51 event.
type Event struct {
	reg *mmio.U32
}

// GetEvent returns n-th te event.
func GetEvent(ph *internal.Pheader, n int) Event {
	return Event{&ph.Events[n]}
}

// Happened returns happened flag (true if event happened).
func (e Event) Happened() bool {
	return e.reg.Load() != 0
}

// Clear clears happened flag for event.
func (e Event) Clear() {
	e.reg.Store(0)
}

func (e Event) phmask() (*internal.Pheader, uint32) {
	ea := e.reg.Addr()
	pha := ea & 0xfffff000
	n := (ea-pha)>>2 - 64
	return (*internal.Pheader)(unsafe.Pointer(pha)), 1 << n
}

// IntEnabled tells whether the event generates interrupts.
func (e Event) IntEnabled() bool {
	ph, mask := e.phmask()
	return ph.IntEnSet.Load()&mask != 0
}

// EnableInt enables generating interrupts by event.
func (e Event) EnableInt() {
	ph, mask := e.phmask()
	ph.IntEnSet.Store(mask)
}

// DisableInt disables generating interrupts by event.
func (e Event) DisableInt() {
	ph, mask := e.phmask()
	ph.IntEnClr.Store(mask)
}

// IRQ returns IRQ number associated to event
func (e Event) IRQ() nvic.IRQ {
	ph, _ := e.phmask()
	return ph.IRQ()
}

// Enabled tells whether the event is enabled, that is, it can update happened
// flag and will be routed to PPI. Only some peripherals (eg. RTC) implements
// this method.
func (e Event) Enabled() bool {
	ph, mask := e.phmask()
	return ph.EvtEnSet.Load()&mask != 0
}

// Enable enables event as described in Enabled description. Only some
// peripherals (eg. RTC) implements this method. Note that EnableInt also
// enables update of happened flag but doesn't connect event to PPI.
func (e Event) Enable() {
	ph, mask := e.phmask()
	ph.EvtEnSet.Store(mask)
}

// Disable disables event. Only some peripherals (eg. RTC) implements this
// method. Note that Disable doesn't affect interrupt generation if it is
// enabled using EnableInt.
func (e Event) Disable() {
	ph, mask := e.phmask()
	ph.EvtEnClr.Store(mask)
}
