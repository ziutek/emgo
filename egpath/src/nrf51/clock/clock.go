// Package clock provides interface to manage nRF51 clocks source/generation.
package clock

import (
	"arch/cortexm/exce"
	"mmio"
	"unsafe"

	"nrf51/internal"
	"nrf51/ppi"
)

// Periph represents clock management peripheral.
type Periph struct {
	te           internal.TasksEvents
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

type Task byte

const (
	HFCLKSTART Task = 0 // Start high frequency crystal oscilator.
	HFCLKSTOP  Task = 1 // Stop high frequency crystal oscilator.
	LFCLKSTART Task = 2 // Start low frequency source.
	LFCLKSTOP  Task = 3 // Stop low frequency source.
	CAL        Task = 4 // Start calibration of low freq. RC oscilator.
	CTSTART    Task = 5 // Start calibration timer.
	CTSTOP     Task = 6 // Stop calibration timer.
)

type Event byte

const (
	HFCLKSTARTED Event = 0 // High frequency crystal oscilator started.
	LFCLKSTARTED Event = 1 // Low frequency source started.
	DONE         Event = 3 // Calibration of low freq. RC osc. complete.
	CTTO         Event = 4 // Calibration timer timeout.
)

func (p *Periph) Task(n Task) ppi.Task {
	return ppi.GetTask(&p.te, int(n))
}

func (p *Periph) Event(n Event) ppi.Event {
	return ppi.GetEvent(&p.te, int(n))
}

func (p *Periph) IRQ() exce.Exce {
	return p.te.IRQ()
}

// HFCLKRun returns true if HFCLKSTART task was triggered.
func (p *Periph) HFCLKRun() bool {
	return p.hfclkrun.Load() != 0
}

type Src byte

const (
	RC Src = iota
	Xtal
	Synth
)

// HFCLKStat returns information about HFCLK status (running or not) and clock
// source.
func (p *Periph) HFCLKStat() (src Src, running bool) {
	s := p.hfclkstat.Load()
	return Src(s & 1), s&(1<<16) != 0
}

// LFCLKRun returns true if LFCLKSTART task was triggered.
func (p *Periph) LFCLKRun() bool {
	return p.lfclkrun.Bit(0)
}

// LFCLKStat returns information about LFCLK status (running or not) and clock
// source.
func (p *Periph) LFCLKStat() (src Src, running bool) {
	s := p.lfclkstat.Load()
	return Src(s & 1), s&(1<<16) != 0
}

// LFCLKSrcCopy returns clock source for LFCLK from time when LFCLKSTART task
// has been triggered.
func (p *Periph) LFCLKSrcCopy() Src {
	return Src(p.lfclksrccopy.LoadBits(3))
}

// LFCLKSrc returns clock source for LFCLK..
func (p *Periph) LFCLKSrc() Src {
	return Src(p.lfclksrc.LoadBits(3))
}

// SetLFCLKSrc sets clock source for LFCLK. It can only be modified when LFCLK
// is not running.
func (p *Periph) SetLFCLKSrc(src Src) {
	p.lfclksrc.Store(uint32(src))
}

// CTIV returns calibration timer interval in milliseconds.
func (p *Periph) CTIV() int {
	return int(p.ctiv.LoadBits(0x7f) * 250)
}

// SetCTIV sets calibration timer interval as number of milliseconds
// (range: 250 ms to 31750 ms).
func (p *Periph) SetCTIV(ctiv int) {
	p.ctiv.Store(uint32(ctiv+125) / 250)
}

type XtalFreq byte

const (
	XF16MHz XtalFreq = 0xff
	XF32MHz XtalFreq = 0x00
)

func (p *Periph) XtalFreq() XtalFreq {
	return XtalFreq(p.xtalfreq.LoadBits(0xff))
}

func (p *Periph) SetXtalFreq(xf XtalFreq) {
	p.xtalfreq.Store(uint32(xf))
}

var Mgmt = (*Periph)(unsafe.Pointer(internal.BaseAPB + 0x00000))
