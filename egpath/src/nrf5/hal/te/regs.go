package te

import (
	"mmio"
	"unsafe"

	"arch/cortexm/nvic"

	"nrf5/hal/internal/mmap"
)

// Regs should be the first field on any Periph struct.
// It takes 0x400 bytes of memory.
type Regs struct {
	tasks    [32]Task
	_        [32]mmio.U32
	events   [32]Event
	_        [32]mmio.U32
	shorts   mmio.U32
	_        [64]mmio.U32
	intEnSet mmio.U32
	intEnClr mmio.U32
	_        [14]mmio.U32
	evtEnSet mmio.U32
	evtEnClr mmio.U32
	_        [45]mmio.U32
}

func (r *Regs) Task(n int) *Task { return &r.tasks[n] }

func (r *Regs) Event(n int) *Event { return &r.events[n] }

// IRQ returns IRQ number associated to events.
func (r *Regs) IRQ() nvic.IRQ {
	addr := uintptr(unsafe.Pointer(r))
	return nvic.IRQ((addr - mmap.BaseAPB) >> 12)
}

// EventMask is a bitmask that can be used to perform some operation (like
// enable/disable IRQ) on multiple events atomically.
type EventMask uint32

// IRQEnabled returns EventMask, wherein the bit set indicates that the
// corresponding event will generate IRQ.
func (r *Regs) IRQEnabled() EventMask {
	return EventMask(r.intEnSet.Load())
}

// EnableIRQ enables generating of IRQ by events specified by mask.
func (r *Regs) EnableIRQ(mask EventMask) {
	r.intEnSet.Store(uint32(mask))
}

// DisableIRQ disables generating of IRQ by events specified by mask.
func (r *Regs) DisableIRQ(mask EventMask) {
	r.intEnClr.Store(uint32(mask))
}

// PPIEnabled returns EventMask, wherein the bit set indicates that the
// corresponding event will be recorded and will be routed to the PPI.
// Only some peripherals (eg. RTC) implements this method.
func (r *Regs) PPIEnabled() EventMask {
	return EventMask(r.evtEnSet.Load())
}

// EnablePPI enables recording of events specified by mask and routing them to
// the PPI. Only some peripherals (eg. RTC) implement this method. Note that if
// this method is implemented then EnableIRQ method also enables recording of
// events but does not enable routing them to the PPI.
func (r *Regs) EnablePPI(mask EventMask) {
	r.evtEnSet.Store(uint32(mask))
}

// DisablePPI disables recording of events specified by mask and routing them to
// the PPI. Only some peripherals (eg. RTC) implement this method. Note that if
// this method is implemented it doesn't affect interrupt generation enabled by
// EnableIRQ which also enables events recording.
func (r *Regs) DisablePPI(mask EventMask) {
	r.evtEnClr.Store(uint32(mask))
}

// SHORTS returns currents value of the SHORTS register.
func (r *Regs) SHORTS() uint32 {
	return r.shorts.Load()
}

// SetSHORTS stores s to the SHORTS register.
func (r *Regs) SetSHORTS(s uint32) {
	r.shorts.Store(s)
}
