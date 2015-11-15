package timer

import (
	"unsafe"

	"nrf51/periph"
)

// Timer represents timer/counter peripheral.
type Periph struct {
	periph.TasksEvents
	_         [65]uint32
	mode      uint32 // Timer mode selection.
	bitmode   uint32 // Configure the number of bits used by the TIMER.
	_         uint32
	prescaler uint32 // Timer prescaler register.
	_         [11]uint32
	cc        [4]uint32 // Capture/Compare registers.
} //c:volatile

// Tasks
const (
	START    periph.Task = 0  // Star Timer.
	STOP     periph.Task = 1  // Stop Timer.
	COUNT    periph.Task = 2  // Increment Timer (Counter mode only).
	CLEAR    periph.Task = 3  // Clear Timer.
	CAPTURE0 periph.Task = 16 // Capture Timer value to CC[0] register.
	CAPTURE1 periph.Task = 17 // Capture Timer value to CC[1] register.
	CAPTURE2 periph.Task = 18 // Capture Timer value to CC[2] register.
	CAPTURE3 periph.Task = 19 // Capture Timer value to CC[3] register.
)

// Events
const (
	COMPARE0 periph.Event = 16 // Compare event on CC[0] match.
	COMPARE1 periph.Event = 17 // Compare event on CC[1] match.
	COMPARE2 periph.Event = 18 // Compare event on CC[2] match.
	COMPARE3 periph.Event = 19 // Compare event on CC[3] match.
)

// Shorts
const (
	COMPARE0_CLEAR periph.Shorts = 1 << 0
	COMPARE1_CLEAR periph.Shorts = 1 << 1
	COMPARE2_CLEAR periph.Shorts = 1 << 2
	COMPARE3_CLEAR periph.Shorts = 1 << 3
	COMPARE0_STOP  periph.Shorts = 0x100 << 0
	COMPARE1_STOP  periph.Shorts = 0x100 << 1
	COMPARE2_STOP  periph.Shorts = 0x100 << 2
	COMPARE3_STOP  periph.Shorts = 0x100 << 3
)

var (
	Timer0 = (*Periph)(unsafe.Pointer(periph.BaseAPB + 0x08000))
	Timer1 = (*Periph)(unsafe.Pointer(periph.BaseAPB + 0x09000))
	Timer2 = (*Periph)(unsafe.Pointer(periph.BaseAPB + 0x0a000))
)

type Mode byte

const (
	TIMER   Mode = 0
	COUNTER Mode = 1
)

func (p *Periph) Mode() Mode {
	return Mode(p.mode)
}

func (p *Periph) SetMode(m Mode) {
	p.mode = uint32(m)
}

type Bitmode byte

const (
	BIT8  Bitmode = 1
	BIT16 Bitmode = 0
	BIT24 Bitmode = 2
	BIT32 Bitmode = 3
)

func (p *Periph) Bitmode() Bitmode {
	return Bitmode(p.bitmode)
}

func (p *Periph) SetBitmode(m Bitmode) {
	p.bitmode = uint32(m)
}

func (p *Periph) Prescaler() int {
	return int(p.prescaler)
}

// SetPrescaler sets prescaler to exp (freq = 16MHz/2^exp). Must only be used
// when the timer is stopped.
func (p *Periph) SetPrescaler(exp int) {
	p.prescaler = uint32(exp)
}

func (p *Periph) CC(n int) uint32 {
	return p.cc[n]
}

func (p *Periph) SetCC(n int, cc uint32) {
	p.cc[n] = cc
}
