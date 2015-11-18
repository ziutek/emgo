package rtc

import (
	"arch/cortexm/exce"
	"mmio"
	"unsafe"

	"nrf51/internal"
	"nrf51/ppi"
)

// Periph represents Real Time Counter peripheral.
type Periph struct {
	te        internal.TasksEvents
	_         [65]mmio.U32
	counter   mmio.U32 // Current COUNTER value.
	prescaler mmio.U32 // 12 bit prescaler for COUNTER frequency.
	_         [13]mmio.U32
	cc        [4]mmio.U32 // Compare registers.
}

// Tasks

func (p *Periph) START() ppi.Task      { return ppi.GetTask(&p.te, 0) }
func (p *Periph) STOP() ppi.Task       { return ppi.GetTask(&p.te, 1) }
func (p *Periph) CLEAR() ppi.Task      { return ppi.GetTask(&p.te, 2) }
func (p *Periph) TRIGOVRFLW() ppi.Task { return ppi.GetTask(&p.te, 3) }

// Events

func (p *Periph) TICK() ppi.Event         { return ppi.GetEvent(&p.te, 0) }
func (p *Periph) OVRFLW() ppi.Event       { return ppi.GetEvent(&p.te, 1) }
func (p *Periph) COMPARE(n int) ppi.Event { return ppi.GetEvent(&p.te, 16+n) }

func (p *Periph) IRQ() exce.Exce {
	return p.te.IRQ()
}

// Counter returns value of counter register.
func (p *Periph) Counter() uint32 {
	return p.counter.LoadBits(0xffffff)
}

func (p *Periph) SetCounter(c uint32) {
	p.counter.Store(c)
}

// Prescaler returns value of prescaler register.
func (p *Periph) Prescaler() uint32 {
	return p.counter.LoadBits(0xfff)
}

// SetPrescaler sets prescaler to pr (freq = 32768Hz/(pr+1)). Must only be used
// when the timer is stopped.
func (p *Periph) SetPrescaler(pr int) {
	p.prescaler.Store(uint32(pr))
}

func (p *Periph) CC(n int) uint32 {
	return p.cc[n].LoadBits(0xffffff)
}

func (p *Periph) SetCC(n int, cc uint32) {
	p.cc[n].Store(cc)
}

var (
	RTC0 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x0b000))
	RTC1 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x11000))
	RTC2 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x24000))
)
