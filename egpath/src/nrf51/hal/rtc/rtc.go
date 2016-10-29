// Package rtc provides interface to manage nRF51 real time counters.
package rtc

import (
	"arch/cortexm/nvic"
	"mmio"
	"unsafe"

	"nrf51/hal/internal"
	"nrf51/hal/te"
)

// Periph represents Real Time Counter peripheral.
type Periph struct {
	ph        internal.Pheader
	_         [65]mmio.U32
	counter   mmio.U32
	prescaler mmio.U32
	_         [13]mmio.U32
	cc        [4]mmio.U32
}

type Task int

const (
	START      Task = 0 // Start RTC COUNTER.
	STOP       Task = 1 // Stop RTC COUNTER.
	CLEAR      Task = 2 // Clear RTC COUNTER.
	TRIGOVRFLW Task = 3 // Set COUNTER to 0xFFFFF0.
)

type Event int

const (
	TICK     Event = 0  // Event on COUNTER increment.
	OVRFLW   Event = 1  // Event on COUNTER overflow.
	COMPARE0 Event = 16 // Compare event on CC[0] match.
	COMPARE1 Event = 17 // Compare event on CC[1] match.
	COMPARE2 Event = 18 // Compare event on CC[2] match.
	COMPARE3 Event = 19 // Compare event on CC[3] match.
)

func (p *Periph) IRQ() nvic.IRQ              { return p.ph.IRQ() }
func (p *Periph) TASK(n Task) *te.TaskReg    { return te.GetTaskReg(&p.ph, int(n)) }
func (p *Periph) EVENT(n Event) *te.EventReg { return te.GetEventReg(&p.ph, int(n)) }

// COUNTER returns value of counter register.
func (p *Periph) COUNTER() uint32 {
	return p.counter.Bits(0xffffff)
}

// SetCOUNTER sets value of counter register.
func (p *Periph) SetCOUNTER(c uint32) {
	p.counter.Store(c)
}

// PRESCALER returns value of prescaler register.
func (p *Periph) PRESCALER() uint32 {
	return p.counter.Bits(0xfff)
}

// SetPRESCALER sets prescaler to pr (freq = 32768Hz/(pr+1)). Must only be used
// when the timer is stopped.
func (p *Periph) SetPRESCALER(pr int) {
	p.prescaler.Store(uint32(pr))
}

// CC returns value of n-th compare register.
func (p *Periph) CC(n int) uint32 {
	return p.cc[n].Bits(0xffffff)
}

// SetCC sets n-th compare register to cc.
func (p *Periph) SetCC(n int, cc uint32) {
	p.cc[n].Store(cc)
}

var (
	RTC0 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x0b000))
	RTC1 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x11000))
	RTC2 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x24000))
)
