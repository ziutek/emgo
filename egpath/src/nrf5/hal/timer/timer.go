// Package timer provides interface to manage nRF5 timers.
package timer

import (
	"mmio"
	"unsafe"

	"nrf5/hal/internal"
	"nrf5/hal/te"
)

// Periph represents timer/counter peripheral.
type Periph struct {
	te.Regs

	_         [65]mmio.U32
	mode      mmio.U32 // Timer mode selection.
	bitmode   mmio.U32 // Configure the number of bits used by the TIMER.
	_         mmio.U32
	prescaler mmio.U32 // Timer prescaler register.
	_         [11]mmio.U32
	cc        [4]mmio.U32 // Capture/Compare registers.
}

//emgo:const
var (
	TIMER0 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x08000))
	TIMER1 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x09000))
	TIMER2 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x0a000))
)

type Task int

const (
	START    Task = 0  // Start Timer.
	STOP     Task = 1  // Stop Timer.
	COUNT    Task = 2  // Increment Timer (Counter mode only).
	CLEAR    Task = 3  // Clear timer.
	CAPTURE0 Task = 16 // Capture Timer value to CC0 register.
	CAPTURE1 Task = 17 // Capture Timer value to CC1 register.
	CAPTURE2 Task = 18 // Capture Timer value to CC2 register.
	CAPTURE3 Task = 19 // Capture Timer value to CC3 register.
)

type Event int

const (
	COMPARE0 Event = 16 // Compare event on CC[0] match.
	COMPARE1 Event = 17 // Compare event on CC[1] match.
	COMPARE2 Event = 18 // Compare event on CC[2] match.
	COMPARE3 Event = 19 // Compare event on CC[3] match.
)

func (p *Periph) Task(t Task) *te.Task      { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event   { return p.Regs.Event(int(e)) }

type Shorts uint32

const (
	COMPARE0_CLEAR Shorts = 1 << 0
	COMPARE1_CLEAR Shorts = 1 << 1
	COMPARE2_CLEAR Shorts = 1 << 2
	COMPARE3_CLEAR Shorts = 1 << 3
	COMPARE0_STOP  Shorts = 1 << 8
	COMPARE1_STOP  Shorts = 1 << 9
	COMPARE2_STOP  Shorts = 1 << 10
	COMPARE3_STOP  Shorts = 1 << 11
)

func (p *Periph) SHORTS() Shorts     { return Shorts(p.Regs.SHORTS()) }
func (p *Periph) SetSHORTS(s Shorts) { p.Regs.SetSHORTS(uint32(s)) }

type Mode byte

const (
	TIMER   Mode = 0
	COUNTER Mode = 1
)

func (p *Periph) MODE() Mode {
	return Mode(p.mode.Bits(1))
}

func (p *Periph) SetMODE(m Mode) {
	p.mode.Store(uint32(m))
}

type Bitmode byte

const (
	BIT8  Bitmode = 1
	BIT16 Bitmode = 0
	BIT24 Bitmode = 2
	BIT32 Bitmode = 3
)

func (p *Periph) BITMODE() Bitmode {
	return Bitmode(p.bitmode.Bits(3))
}

func (p *Periph) SetBITMODE(m Bitmode) {
	p.bitmode.Store(uint32(m))
}

func (p *Periph) PRESCALER() int {
	return int(p.prescaler.Bits(0xf))
}

// SetPrescaler sets prescaler to exp (freq = 16MHz/2^exp). Must only be used
// when the timer is stopped.
func (p *Periph) SetPRESCALER(exp int) {
	p.prescaler.Store(uint32(exp))
}

func (p *Periph) CC(n int) uint32 {
	return p.cc[n].Load()
}

func (p *Periph) SetCC(n int, cc uint32) {
	p.cc[n].Store(cc)
}
