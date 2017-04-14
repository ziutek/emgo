// Package timer provides interface to manage nRF5 timers.
package timer

import (
	"mmio"
	"unsafe"

	"nrf5/hal/internal/mmap"
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
	cc        [6]mmio.U32 // Capture/Compare registers.
}

//emgo:const
var (
	TIMER0 = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x08000))
	TIMER1 = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x09000))
	TIMER2 = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x0A000))
	TIMER3 = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x1A000)) // nRF52.
	TIMER4 = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x1B000)) // nRF52.
)

type Task int

const (
	START Task = 0 // Start Timer.
	STOP  Task = 1 // Stop Timer.
	COUNT Task = 2 // Increment Timer (Counter mode only).
	CLEAR Task = 3 // Clear timer.
)

func (p *Periph) Task(t Task) *te.Task { return p.Regs.Task(int(t)) }

// CAPTURE returns Capture task for CCn register.
func (p *Periph) CAPTURE(n int) *te.Task {
	return p.Regs.Task(16 + n)
}

// COMPARE returns Compare event for CCn register.
func (p *Periph) COMPARE(n int) *te.Event {
	return p.Regs.Event(16 + n)
}

type Shorts uint32

const (
	COMPARE0_CLEAR Shorts = 1 << 0
	COMPARE1_CLEAR Shorts = 1 << 1
	COMPARE2_CLEAR Shorts = 1 << 2
	COMPARE3_CLEAR Shorts = 1 << 3
	COMPARE4_CLEAR Shorts = 1 << 4
	COMPARE5_CLEAR Shorts = 1 << 5
	COMPARE0_STOP  Shorts = 1 << 8
	COMPARE1_STOP  Shorts = 1 << 9
	COMPARE2_STOP  Shorts = 1 << 10
	COMPARE3_STOP  Shorts = 1 << 11
	COMPARE4_STOP  Shorts = 1 << 12
	COMPARE5_STOP  Shorts = 1 << 13
)

func (p *Periph) LoadSHORTS() Shorts   { return Shorts(p.Regs.LoadSHORTS()) }
func (p *Periph) StoreSHORTS(s Shorts) { p.Regs.StoreSHORTS(uint32(s)) }

type Mode byte

const (
	TIMER   Mode = 0
	COUNTER Mode = 1
)

func (p *Periph) LoadMODE() Mode {
	return Mode(p.mode.Bits(1))
}

func (p *Periph) StoreMODE(m Mode) {
	p.mode.Store(uint32(m))
}

type Bitmode byte

const (
	BIT8  Bitmode = 1
	BIT16 Bitmode = 0
	BIT24 Bitmode = 2
	BIT32 Bitmode = 3
)

func (p *Periph) LoadBITMODE() Bitmode {
	return Bitmode(p.bitmode.Bits(3))
}

func (p *Periph) StoreBITMODE(m Bitmode) {
	p.bitmode.Store(uint32(m))
}

func (p *Periph) LoadPRESCALER() int {
	return int(p.prescaler.Bits(0xf))
}

// StorePRESCALER sets prescaler to exp (freq = 16MHz/2^exp). Must only be used
// when the timer is stopped.
func (p *Periph) StorePRESCALER(exp int) {
	p.prescaler.Store(uint32(exp))
}

// LoadCC returns value of n-th Capture/Compare register. nRF51/nRF52 has 4/6 CC
// registers
func (p *Periph) LoadCC(n int) uint32 {
	return p.cc[n].Load()
}

// StoreCC stores cc into n-th Capture/Compare register. nRF51/nRF52 has 4/6 CC
// registers
func (p *Periph) StoreCC(n int, cc uint32) {
	p.cc[n].Store(cc)
}
