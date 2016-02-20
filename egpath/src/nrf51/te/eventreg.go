package te

import (
	"mmio"
	"unsafe"

	"arch/cortexm/nvic"

	"nrf51/internal"
)

// EventReg represents peripheral registers that are used to records an
// occurences of some kind of events.
type EventReg struct {
	u32 mmio.U32
}

// GetEventReg is for internal use.
func GetEventReg(ph *internal.Pheader, n int) *EventReg {
	return (*EventReg)(unsafe.Pointer(&ph.Events[n]))
}

// IsSet reports whether r recorded occurrence of an event.
func (r *EventReg) IsSet() bool {
	return r.u32.Load() != 0
}

// Clear clears r.
func (r *EventReg) Clear() {
	r.u32.Store(0)
}

func phmask(r *EventReg) (*internal.Pheader, uint32) {
	ea := r.u32.Addr()
	pha := ea & 0xfffff000
	n := (ea-pha)>>2 - 64
	return (*internal.Pheader)(unsafe.Pointer(pha)), 1 << n
}

// IntEnabled reports whether the occurence of an event will generat interrupt.
func (r *EventReg) IntEnabled() bool {
	ph, mask := phmask(r)
	return ph.IntEnSet.Load()&mask != 0
}

// EnableInt enables generating interrupts by event recorded by r.
func (r *EventReg) EnableInt() {
	ph, mask := phmask(r)
	ph.IntEnSet.Store(mask)
}

// DisableInt disables generating interrupts by event recorded by r.
func (r *EventReg) DisableInt() {
	ph, mask := phmask(r)
	ph.IntEnClr.Store(mask)
}

// IRQ returns IRQ number associated to r.
func (r *EventReg) IRQ() nvic.IRQ {
	ph, _ := phmask(r)
	return ph.IRQ()
}

// EventEnabled reports whether the occurrence of an event will be recorded by
// r and will be routed to PPI. Only some peripherals (eg. RTC) implements
// this method.
func (r *EventReg) EventEnabled() bool {
	ph, mask := phmask(r)
	return ph.EvtEnSet.Load()&mask != 0
}

// EnableEvent enables recording of events and routing them to PPI. Only some
// peripherals (eg. RTC) implements this method. Note that if this method is
// implemented then EnableInt method also enables recording ot events but does
// not enable routing them to PPI.
func (r *EventReg) EnableEvent() {
	ph, mask := phmask(r)
	ph.EvtEnSet.Store(mask)
}

// DisableEvent disables recording of events and routing them to PPI. Only some
// peripherals (eg. RTC) implements this method. Note that if this method is
// implemented it doesn't affect interrupt generation enabled by EnableInt which
// also enables events recording.
func (r *EventReg) DisableEvent() {
	ph, mask := phmask(r)
	ph.EvtEnClr.Store(mask)
}
