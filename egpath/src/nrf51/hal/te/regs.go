package te

import (
	"mmio"
	"unsafe"

	"arch/cortexm/nvic"

	"nrf51/hal/internal"
)

// Regs should be the first field on any Periph struct.
// It takes 0x400 bytes of memory.
type Regs struct {
	tasks    [32]TaskReg
	_        [32]mmio.U32
	events   [32]EventReg
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

func (r *Regs) TASK(n int) *TaskReg { return &r.tasks[n] }

func (r *Regs) EVENT(n int) *EventReg { return &r.events[n] }

// IRQ returns IRQ number associated to events.
func (r *Regs) IRQ() nvic.IRQ {
	addr := uintptr(unsafe.Pointer(r))
	return nvic.IRQ((addr - internal.BaseAPB) >> 12)
}

/*
type ShortReg struct{ u32 mmio.U32 }

func (r *ShortsReg) Load() Shorts   { return Shorts(r.u32.Load()) }
func (r *ShortsReg) Store(s Shorts) { r.u32.Store(uint32(s)) }
func (r *ShortsReg) Set(s Shorts)   { r.u32.SetBits(uint32(s)) }
func (r *ShortsReg) Clear(s Shorts) { r.u32.ClearBits(uint32(s)) }
*/