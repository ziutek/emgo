// Package clock provides interface to manage nRF51 clocks source/generation.
package clock

import (
	"unsafe"

	"nrf51/periph"
)

// Periph represents clock management peripheral.
type Periph struct {
	periph.TasksEvents
	_            [2]uint32
	hfclkrun     uint32
	hfclkstat    uint32
	_            uint32
	lfclkrun     uint32
	lfclkstat    uint32
	lfclksrccopy uint32
	_            [62]uint32
	lfclksrc     uint32
	_            [7]uint32
	ctiv         uint32
	_            [5]uint32
	xtalfreq     uint32
} //c:volatile

// Tasks
const (
	HFCLKSTART periph.Task = 0 // Start high frequency crystal oscilator.
	HFCLKSTOP  periph.Task = 1 // Stop high frequency crystal oscilator.
	LFCLKSTART periph.Task = 2 // Start low frequency source.
	LFCLKSTOP  periph.Task = 3 // Stop low frequency source.
	CAL        periph.Task = 4 // Start calibration of low freq. RC oscilator.
	CTSTART    periph.Task = 5 // Start calibration timer.
	CTSTOP     periph.Task = 6 // Stop calibration timer.
)

// Events
const (
	HFCLKSTARTED periph.Event = 0 // High frequency crystal oscilator started.
	LFCLKSTARTED periph.Event = 1 // Low frequency source started.
	DONE         periph.Event = 3 // Calibration of low freq. RC osc. complete.
	CTTO         periph.Event = 4 // Calibration timer timeout.
)

// HFCLKRun returns true if HFCLKSTART task was triggered.
func (p *Periph) HFCLKRun() bool {
	return p.hfclkrun != 0
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
	s := p.hfclkstat
	return Src(s & 1), s&(1<<16) != 0
}

// LFCLKRun returns true if LFCLKSTART task was triggered.
func (p *Periph) LFCLKRun() bool {
	return p.lfclkrun != 0
}

// LFCLKStat returns information about LFCLK status (running or not) and clock
// source.
func (p *Periph) LFCLKStat() (src Src, running bool) {
	s := p.lfclkstat
	return Src(s & 1), s&(1<<16) != 0
}

// LFCLKSrcCopy returns clock source for LFCLK from time when LFCLKSTART task
// has been triggered.
func (p *Periph) LFCLKSrcCopy() Src {
	return Src(p.lfclksrccopy & 3)
}

// LFCLKSrc returns clock source for LFCLK..
func (p *Periph) LFCLKSrc() Src {
	return Src(p.lfclksrc & 3)
}

// SetLFCLKSrc sets clock source for LFCLK. It can only be modified when LFCLK
// is not running.
func (p *Periph) SetLFCLKSrc(src Src) {
	p.lfclksrc = uint32(src)
}

// CTIV returns calibration timer interval in milliseconds.
func (p *Periph) CTIV() int {
	return int(p.ctiv * 250)
}

// SetCTIV sets calibration timer interval as number of milliseconds
// (range: 250 ms to 31750 ms).
func (p *Periph) SetCTIV(ctiv int) {
	p.ctiv = uint32(ctiv+125) / 250
}

type XtalFreq byte

const (
	XF16MHz XtalFreq = 0xff
	XF32MHz XtalFreq = 0x00
)

func (p *Periph) XtalFreq() XtalFreq {
	return XtalFreq(p.xtalfreq)
}

func (p *Periph) SetXtalFreq(xf XtalFreq) {
	p.xtalfreq = uint32(xf)
}

var Mgmt = (*Periph)(unsafe.Pointer(periph.BaseAPB + 0x00000))
