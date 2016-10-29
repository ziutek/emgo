// Package clock provides interface to manage nRF51 clocks source/generation.
package clock

import (
	"arch/cortexm/nvic"
	"mmio"
	"unsafe"

	"nrf51/hal/internal"
	"nrf51/hal/te"
)

// Periph represents clock management peripheral.
type Periph struct {
	ph           internal.Pheader
	_            [2]mmio.U32
	hfclkrun     mmio.U32
	hfclkstat    mmio.U32
	_            mmio.U32
	lfclkrun     mmio.U32
	lfclkstat    mmio.U32
	lfclksrccopy mmio.U32
	_            [62]mmio.U32
	lfclksrc     mmio.U32
	_            [7]mmio.U32
	ctiv         mmio.U32
	_            [5]mmio.U32
	xtalfreq     mmio.U32
}

type Task int

const (
	HFCLKSTART Task = 0 // Start high frequency crystal oscilator.
	HFCLKSTOP  Task = 1 // Stop high frequency crystal oscilator.
	LFCLKSTART Task = 2 // Start low frequency source.
	LFCLKSTOP  Task = 3 // Stop low frequency source.
	CAL        Task = 4 // Start calibration of low freq. RC oscilator.
	CTSTART    Task = 5 // Start calibration timer.
	CTSTOP     Task = 6 // Stop calibration timer.
)

type Event int

const (
	HFCLKSTARTED Event = 0 // High frequency crystal oscilator started.
	LFCLKSTARTED Event = 1 // Low frequency source started.
	DONE         Event = 3 // Calibration of low freq. RC osc. complete.
	CTTO         Event = 4 // Calibration timer timeout.
)

func (p *Periph) IRQ() nvic.IRQ              { return p.ph.IRQ() }
func (p *Periph) TASK(n Task) *te.TaskReg    { return te.GetTaskReg(&p.ph, int(n)) }
func (p *Periph) EVENT(n Event) *te.EventReg { return te.GetEventReg(&p.ph, int(n)) }

// HFCLKRUN returns true if HFCLKSTART task was triggered.
func (p *Periph) HFCLKRUN() bool {
	return p.hfclkrun.Load() != 0
}

type SRC byte

const (
	RC    SRC = 0
	Xtal  SRC = 1
	Synth SRC = 2
)

// HFCLKStat returns information about HFCLK status (running or not) and clock
// source.
func (p *Periph) HFCLKSTAT() (src SRC, running bool) {
	s := p.hfclkstat.Load()
	return SRC(s & 1), s&(1<<16) != 0
}

// LFCLKRUN returns true if LFCLKSTART task was triggered.
func (p *Periph) LFCLKRUN() bool {
	return p.lfclkrun.Bit(0) != 0
}

// LFCLKSTAT returns information about LFCLK status (running or not) and clock
// source.
func (p *Periph) LFCLKSTAT() (src SRC, running bool) {
	s := p.lfclkstat.Load()
	return SRC(s & 1), s&(1<<16) != 0
}

// LFCLKSRCCOPY returns clock source for LFCLK from time when LFCLKSTART task
// has been triggered.
func (p *Periph) LFCLKSRCCOPY() SRC {
	return SRC(p.lfclksrccopy.Bits(3))
}

// LFCLKSRC returns clock source for LFCLK..
func (p *Periph) LFCLKSRC() SRC {
	return SRC(p.lfclksrc.Bits(3))
}

// SetLFCLKSRC sets clock source for LFCLK. It can only be modified when LFCLK
// is not running.
func (p *Periph) SetLFCLKSRC(src SRC) {
	p.lfclksrc.Store(uint32(src))
}

// CTIV returns calibration timer interval in milliseconds.
func (p *Periph) CTIV() int {
	return int(p.ctiv.Bits(0x7f) * 250)
}

// SetCTIV sets calibration timer interval as number of milliseconds
// (range: 250 ms to 31750 ms).
func (p *Periph) SetCTIV(ctiv int) {
	p.ctiv.Store(uint32(ctiv+125) / 250)
}

type XTALFREQ byte

const (
	XF16MHz XTALFREQ = 0xff
	XF32MHz XTALFREQ = 0x00
)

func (p *Periph) XTALFREQ() XTALFREQ {
	return XTALFREQ(p.xtalfreq.Bits(0xff))
}

func (p *Periph) SetXTALFREQ(xf XTALFREQ) {
	p.xtalfreq.Store(uint32(xf))
}

var Mgmt = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x00000))
