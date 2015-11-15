package rtc

import (
	"unsafe"

	"nrf51/periph"
)

// Periph represents Real Time Counter peripheral.
type Periph struct {
	periph.TasksEvents
	_         [65]uint32
	counter   uint32 // Current COUNTER value.
	prescaler uint32 // 12 bit prescaler for COUNTER frequency.
	_         [13]uint32
	cc        [4]uint32 // Compare registers.
} //c:volatile

// Tasks
const (
	START      periph.Task = 0 // Start RTC COUNTER.
	STOP       periph.Task = 1 // Stop RTC COUNTER.
	CLEAR      periph.Task = 2 // Clear RTC COUNTER.
	TRIGOVRFLW periph.Task = 3 // Set COUNTER to 0xFFFFF0.
)

// Events
const (
	TICK     periph.Event = 0  // Event on COUNTER increment.
	OVRFLW   periph.Event = 1  // Event on COUNTER overflow.
	COMPARE0 periph.Event = 16 // Compare event on CC[0] match.
	COMPARE1 periph.Event = 17 // Compare event on CC[1] match.
	COMPARE2 periph.Event = 18 // Compare event on CC[2] match.
	COMPARE3 periph.Event = 19 // Compare event on CC[3] match.
)

var (
	RTC0 = (*Periph)(unsafe.Pointer(periph.BaseAPB + 0x0b000))
	RTC1 = (*Periph)(unsafe.Pointer(periph.BaseAPB + 0x11000))
	RTC2 = (*Periph)(unsafe.Pointer(periph.BaseAPB + 0x24000))
)

func (p *Periph) Counter() uint32 {
	return p.counter
}

func (p *Periph) SetCounter(c uint32) {
	p.counter = c
}

// SetPrescaler sets prescaler to pr (freq = 32768Hz/(pr+1)). Must only be used
// when the timer is stopped.
func (p *Periph) SetPrescaler(pr int) {
	p.prescaler = uint32(pr)
}

func (p *Periph) CC(n int) uint32 {
	return p.cc[n]
}

func (p *Periph) SetCC(n int, cc uint32) {
	p.cc[n] = cc
}
