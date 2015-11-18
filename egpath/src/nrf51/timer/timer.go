package timer

import (
	"arch/cortexm/exce"
	"mmio"
	"unsafe"

	"nrf51/internal"
	"nrf51/ppi"
)

// Timer represents timer/counter peripheral.
type Periph struct {
	te        internal.TasksEvents
	_         [65]mmio.U32
	mode      mmio.U32 // Timer mode selection.
	bitmode   mmio.U32 // Configure the number of bits used by the TIMER.
	_         mmio.U32
	prescaler mmio.U32 // Timer prescaler register.
	_         [11]mmio.U32
	cc        [4]mmio.U32 // Capture/Compare registers.
}

// Tasks

func (p *Periph) START() ppi.Task        { return ppi.GetTask(&p.te, 0) }
func (p *Periph) STOP() ppi.Task         { return ppi.GetTask(&p.te, 1) }
func (p *Periph) COUNT() ppi.Task        { return ppi.GetTask(&p.te, 2) }
func (p *Periph) CLEAR() ppi.Task        { return ppi.GetTask(&p.te, 3) }
func (p *Periph) CAPTURE(n int) ppi.Task { return ppi.GetTask(&p.te, 16+n) }

// Events

func (p *Periph) COMPARE(n int) ppi.Event { return ppi.GetEvent(&p.te, 16+n) }

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

func (p *Periph) Shorts() Shorts {
	return Shorts(p.te.Shorts.Load())
}

func (p *Periph) SetShorts(s Shorts) {
	p.te.Shorts.SetBits(uint32(s))
}

func (p *Periph) ClearShorts(s Shorts) {
	p.te.Shorts.ClearBits(uint32(s))
}

func (p *Periph) IRQ() exce.Exce {
	return p.te.IRQ()
}

type Mode byte

const (
	TIMER   Mode = 0
	COUNTER Mode = 1
)

func (p *Periph) Mode() Mode {
	return Mode(p.mode.LoadBits(1))
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
	return Bitmode(p.bitmode.LoadBits(3))
}

func (p *Periph) SetBitmode(m Bitmode) {
	p.bitmode.Store(uint32(m))
}

func (p *Periph) Prescaler() int {
	return int(p.prescaler.LoadBits(0xf))
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
