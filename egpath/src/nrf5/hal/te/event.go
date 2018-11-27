package te

import (
	"mmio"
	"unsafe"

	"arch/cortexm/nvic"
)

// Event represents a peripheral register that is used to records an occurence
// of some kind of event.
type Event struct {
	u32 mmio.U32
}

// IsSet reports whether r recorded occurrence of an event.
func (r *Event) IsSet() bool {
	return r.u32.Load() != 0
}

// Clear clears r.
func (r *Event) Clear() {
	r.u32.Store(0)
}

func regsMask(r *Event) (*Regs, uint32) {
	ea := r.u32.Addr()
	ra := ea & 0xfffff000
	n := (ea-ra)>>2 - 64
	return (*Regs)(unsafe.Pointer(ra)), 1 << n
}

// IntEnabled reports whether the occurence of an event will generate IRQ.
func (r *Event) IRQEnabled() bool {
	rr, mask := regsMask(r)
	return rr.intEnSet.Load()&mask != 0
}

// EnableIRQ enables generating of IRQ by event recorded by r.
func (r *Event) EnableIRQ() {
	rr, mask := regsMask(r)
	rr.intEnSet.Store(mask)
}

// DisableIRQ disables generating of IRQ by event recorded by r.
func (r *Event) DisableIRQ() {
	rr, mask := regsMask(r)
	rr.intEnClr.Store(mask)
}

// NVIRQ returns NVIC IRQ number associated to r.
func (r *Event) NVIRQ() nvic.IRQ {
	rr, _ := regsMask(r)
	return rr.NVIRQ()
}

// PPIEnabled reports whether the occurrence of an event will be recorded by
// r and will be routed to the PPI. Only some peripherals (eg. RTC) implement
// this method.
func (r *Event) PPIEnabled() bool {
	rr, mask := regsMask(r)
	return rr.evtEnSet.Load()&mask != 0
}

// EnablePPI enables recording of events and routing them to the PPI. Only some
// peripherals (eg. RTC) implement this method. Note that if this method is
// implemented then EnableIRQ method also enables recording of events but does
// not enable routing them to the PPI.
func (r *Event) EnablePPI() {
	rr, mask := regsMask(r)
	rr.evtEnSet.Store(mask)
}

// DisablePPI disables recording of events and routing them to the PPI. Only
// some peripherals (eg. RTC) implement this method. Note that if this method
// is implemented it doesn't affect interrupt generation enabled by EnableIRQ
// which also enables events recording.
func (r *Event) DisablePPI() {
	rr, mask := regsMask(r)
	rr.evtEnClr.Store(mask)
}
