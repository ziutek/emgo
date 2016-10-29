package te

import (
	"mmio"
	"unsafe"

	"arch/cortexm/nvic"
)

// EventReg represents peripheral registers that are used to records an
// occurences of some kind of events.
type EventReg struct {
	u32 mmio.U32
}

// IsSet reports whether r recorded occurrence of an event.
func (r *EventReg) IsSet() bool {
	return r.u32.Load() != 0
}

// Clear clears r.
func (r *EventReg) Clear() {
	r.u32.Store(0)
}

func regsMask(r *EventReg) (*Regs, uint32) {
	ea := r.u32.Addr()
	ra := ea & 0xfffff000
	n := (ea-ra)>>2 - 64
	return (*Regs)(unsafe.Pointer(ra)), 1 << n
}

// IntEnabled reports whether the occurence of an event will generate IRQ.
func (r *EventReg) IRQEnabled() bool {
	rr, mask := regsMask(r)
	return rr.intEnSet.Load()&mask != 0
}

// EnableIRQ enables generating of IRQ by event recorded by r.
func (r *EventReg) EnableIRQ() {
	rr, mask := regsMask(r)
	rr.intEnSet.Store(mask)
}

// DisableIRQ disables generating of IRQ by event recorded by r.
func (r *EventReg) DisableIRQ() {
	rr, mask := regsMask(r)
	rr.intEnClr.Store(mask)
}

// IRQ returns IRQ number associated to r.
func (r *EventReg) IRQ() nvic.IRQ {
	rr, _ := regsMask(r)
	return rr.IRQ()
}

// PPIEnabled reports whether the occurrence of an event will be recorded by
// r and will be routed to PPI. Only some peripherals (eg. RTC) implements
// this method.
func (r *EventReg) PPIEnabled() bool {
	rr, mask := regsMask(r)
	return rr.evtEnSet.Load()&mask != 0
}

// Enable enables recording of events and routing them to PPI. Only some
// peripherals (eg. RTC) implements this method. Note that if this method is
// implemented then EnableIRQ method also enables recording of events but does
// not enable routing them to PPI.
func (r *EventReg) EnablePPI() {
	rr, mask := regsMask(r)
	rr.evtEnSet.Store(mask)
}

// Disable disables recording of events and routing them to PPI. Only some
// peripherals (eg. RTC) implements this method. Note that if this method is
// implemented it doesn't affect interrupt generation enabled by EnableIRQ which
// also enables events recording.
func (r *EventReg) DisablePPI() {
	rr, mask := regsMask(r)
	rr.evtEnClr.Store(mask)
}
