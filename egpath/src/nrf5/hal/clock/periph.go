// Package clock provides interface to manage nRF51 clocks source/generation.
package clock

import (
	"mmio"
	"unsafe"

	"nrf5/hal/internal/mmap"
	"nrf5/hal/te"
)

// Periph represents clock management peripheral.
type Periph struct {
	te.Regs

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

//emgo:const
var CLOCK = (*Periph)(unsafe.Pointer(mmap.BaseAPB + 0x00000))

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

func (p *Periph) Task(t Task) *te.Task    { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

// HFCLKRUN returns true if HFCLKSTART task was triggered.
func (p *Periph) LoadHFCLKRUN() bool {
	return p.hfclkrun.Load() != 0
}

type Source byte

const (
	RC    Source = 0
	XTAL  Source = 1
	SYNTH Source = 2
)

// LoadHFCLKStat returns information about HFCLK status (running or not) and
// clock source.
func (p *Periph) LoadHFCLKSTAT() (src Source, running bool) {
	s := p.hfclkstat.Load()
	return Source(s & 1), s&(1<<16) != 0
}

// LoadLFCLKRUN returns true if LFCLKSTART task was triggered.
func (p *Periph) LoadLFCLKRUN() bool {
	return p.lfclkrun.Bit(0) != 0
}

// LoadLFCLKSTAT returns information about LFCLK status (running or not) and
// clock source.
func (p *Periph) LoadLFCLKSTAT() (src Source, running bool) {
	s := p.lfclkstat.Load()
	return Source(s & 1), s&(1<<16) != 0
}

// LoadLFCLKSRCCOPY returns clock source for LFCLK from time when LFCLKSTART
// task has been triggered.
func (p *Periph) LoadLFCLKSRCCOPY() Source {
	return Source(p.lfclksrccopy.Bits(3))
}

// LoadLFCLKSRC returns clock source for LFCLK.
func (p *Periph) LoadLFCLKSRC() Source {
	return Source(p.lfclksrc.Bits(3))
}

// StoreLFCLKSRC sets clock source for LFCLK. It can only be modified when
// LFCLK is not running.
func (p *Periph) StoreLFCLKSRC(src Source) {
	p.lfclksrc.Store(uint32(src))
}

// LoadCTIV returns calibration timer interval in milliseconds.
func (p *Periph) LoadCTIV() int {
	return int(p.ctiv.Bits(0x7f) * 250)
}

// StoreCTIV sets calibration timer interval as number of milliseconds
// (range: 250 ms to 31750 ms).
func (p *Periph) StoreCTIV(ctiv int) {
	p.ctiv.Store(uint32(ctiv+125) / 250)
}

type Freq byte

const (
	F16MHz Freq = 0xff
	F32MHz Freq = 0x00
)

func (p *Periph) LoadXTALFREQ() Freq {
	return Freq(p.xtalfreq.Bits(0xff))
}

func (p *Periph) StoreXTALFREQ(f Freq) {
	p.xtalfreq.Store(uint32(f))
}
