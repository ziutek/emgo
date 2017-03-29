// Package rtc provides interface to nRF5 real time counters.
package rtc

import (
	"mmio"
	"unsafe"

	"nrf5/hal/internal/mmap"
	"nrf5/hal/te"
)

// Periph represents Real Time Counter peripheral.
type Periph struct {
	te.Regs

	_         [65]mmio.U32
	counter   mmio.U32
	prescaler mmio.U32
	_         [13]mmio.U32
	cc        [4]mmio.U32
}

//emgo:const
var (
	RTC0 = (*Periph)(unsafe.Pointer(mmap.BaseAPB + 0x0b000))
	RTC1 = (*Periph)(unsafe.Pointer(mmap.BaseAPB + 0x11000))
	RTC2 = (*Periph)(unsafe.Pointer(mmap.BaseAPB + 0x24000))
)

type Task byte

const (
	START      Task = 0 // Start RTC COUNTER.
	STOP       Task = 1 // Stop RTC COUNTER.
	CLEAR      Task = 2 // Clear RTC COUNTER.
	TRIGOVRFLW Task = 3 // Store COUNTER to 0xFFFFF0.
)

type Event byte

const (
	TICK     Event = 0  // Event on COUNTER increment.
	OVRFLW   Event = 1  // Event on COUNTER overflow.
	COMPARE0 Event = 16 // Compare event on CC[0] match.
	COMPARE1 Event = 17 // Compare event on CC[1] match.
	COMPARE2 Event = 18 // Compare event on CC[2] match.
	COMPARE3 Event = 19 // Compare event on CC[3] match.
)

func (p *Periph) Task(t Task) *te.Task    { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

// LoadCOUNTER returns value of counter register.
func (p *Periph) LoadCOUNTER() uint32 {
	return p.counter.Bits(0xffffff)
}

// StoreCOUNTER stores value of counter register.
func (p *Periph) StoreCOUNTER(c uint32) {
	p.counter.Store(c)
}

// LoadPRESCALER returns value of prescaler register.
func (p *Periph) LoadPRESCALER() uint32 {
	return p.counter.Bits(0xfff)
}

// StorePRESCALER stores prescaler to pr (freq = 32768Hz/(pr+1)). Must only be used
// when the timer is stopped.
func (p *Periph) StorePRESCALER(pr int) {
	p.prescaler.Store(uint32(pr))
}

// LoadCC returns value of n-th compare register.
func (p *Periph) LoadCC(n int) uint32 {
	return p.cc[n].Bits(0xffffff)
}

// StoreCC stores n-th compare register to cc.
func (p *Periph) StoreCC(n int, cc uint32) {
	p.cc[n].Store(cc)
}
