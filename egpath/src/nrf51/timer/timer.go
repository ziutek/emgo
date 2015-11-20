// Package timer provides interface to manage nRF51 timers.
package timer

import (
	"arch/cortexm/exce"
	"mmio"
	"unsafe"

	"nrf51/internal"
	"nrf51/te"
)

// Periph represents timer/counter peripheral.
type Periph struct {
	ph        internal.Pheader
	_         [65]mmio.U32
	mode      mmio.U32 // Timer mode selection.
	bitmode   mmio.U32 // Configure the number of bits used by the TIMER.
	_         mmio.U32
	prescaler mmio.U32 // Timer prescaler register.
	_         [11]mmio.U32
	cc        [4]mmio.U32 // Capture/Compare registers.
}

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

type Shorts uint32

const (
	COMPARE0_CLEAR Shorts = 1 << 0
	COMPARE1_CLEAR Shorts = 1 << 1
	COMPARE2_CLEAR Shorts = 1 << 2
	COMPARE3_CLEAR Shorts = 1 << 3
	COMPARE0_STOP  Shorts = 0x100 << 0
	COMPARE1_STOP  Shorts = 0x100 << 1
	COMPARE2_STOP  Shorts = 0x100 << 2
	COMPARE3_STOP  Shorts = 0x100 << 3
)

func (p *Periph) IRQ() exce.Exce         { return p.ph.IRQ() }
func (p *Periph) Task(n Task) te.Task    { return te.GetTask(&p.ph, int(n)) }
func (p *Periph) Event(n Event) te.Event { return te.GetEvent(&p.ph, int(n)) }
func (p *Periph) Shorts() Shorts         { return Shorts(p.ph.Shorts.Load()) }
func (p *Periph) SetShorts(s Shorts)     { p.ph.Shorts.SetBits(uint32(s)) }
func (p *Periph) ClearShorts(s Shorts)   { p.ph.Shorts.ClearBits(uint32(s)) }

type Mode byte

const (
	TIMER   Mode = 0
	COUNTER Mode = 1
)

func (p *Periph) Mode() Mode {
	return Mode(p.mode.LoadMask(1))
}

func (p *Periph) SetMode(m Mode) {
	p.mode.Store(uint32(m))
}

type Bitmode byte

const (
	BIT8  Bitmode = 1
	BIT16 Bitmode = 0
	BIT24 Bitmode = 2
	BIT32 Bitmode = 3
)

func (p *Periph) Bitmode() Bitmode {
	return Bitmode(p.bitmode.LoadMask(3))
}

func (p *Periph) SetBitmode(m Bitmode) {
	p.bitmode.Store(uint32(m))
}

func (p *Periph) Prescaler() int {
	return int(p.prescaler.LoadMask(0xf))
}

// SetPrescaler sets prescaler to exp (freq = 16MHz/2^exp). Must only be used
// when the timer is stopped.
func (p *Periph) SetPrescaler(exp int) {
	p.prescaler.Store(uint32(exp))
}

func (p *Periph) CC(n int) uint32 {
	return p.cc[n].Load()
}

func (p *Periph) SetCC(n int, cc uint32) {
	p.cc[n].Store(cc)
}

var (
	Timer0 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x08000))
	Timer1 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x09000))
	Timer2 = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x0a000))
)
