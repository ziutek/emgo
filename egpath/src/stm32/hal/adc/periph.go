package adc

import (
	"unsafe"

	"stm32/hal/internal"
	"stm32/hal/system"

	"stm32/hal/raw/adc"
)

type Periph struct {
	raw adc.ADC_Periph
}

// Bus returns a bus to which p is connected to.
func (p *Periph) Bus() system.Bus {
	return internal.Bus(unsafe.Pointer(p))
}

// EnableClock enables clock for p. Lp determines whether the clock remains on
// in low power (sleep) mode. In some MCUs clock cannot be enabled for only one
// ADC, then EnableClock can affects a pair of peripherals (eg.
// ADC1.EnableClock() can affect ADC1 and ADC2 simultaneously).
func (p *Periph) EnableClock(lp bool) {
	p.enableClock(lp)
}

// DisableClock disables clock for p. In some MCUs clock cannot be disabled for
// only one ADC, then DisableClock can affects a pair of peripherals (eg.
// ADC1.DisableClock() can affect ADC1 and ADC2 simultaneously).
func (p *Periph) DisableClock() {
	p.disableClock()
}

// Calibrate calibrates p.
func (p *Periph) Calibrate() {
	p.calibrate()
}

// Enable enables p.
func (p *Periph) Enable() {
	p.enable()
}

// Enabled reports whether p is enabled.
func (p *Periph) Enabled() bool {
	return p.enabled()
}

// Disable disables p. Use Enabled to check that p is effectively disabled.
func (p *Periph) Disable() {
	p.disable()
}

type SamplTime int8

// MaxSamplTime returns the largest possible value of sampling time that takes
// no more than hc half clock cycles. If ht is too small MaxSamplTime returns
// negative value.
func MaxSamplTime(hc int) SamplTime {
	var i int
	for i < len(halfCycles) {
		if int(halfCycles[i]) > hc {
			break
		}
		i++
	}
	return SamplTime(i - 1)
}

// HalfCycles returns number of half clock cycles that st corresponds to.
func (st SamplTime) HalfCycles() int {
	return int(halfCycles[st])
}

// SetSamplTime sets sampling time for channel ch.
func (p *Periph) SetSamplTime(ch int, st SamplTime) {
	p.setSamplTime(ch, st)
}

// SetSequence sets regular sequence of channels.
func (p *Periph) SetSequence(ch ...int) {
	p.setSequence(ch)
}

type TrigSrc byte

// SetTrigSrc selects source of external trigger.
func (p *Periph) SetTrigSrc(src TrigSrc) {
	p.setTrigSrc(src)
}

type TrigEdge byte

const (
	EdgeNone   TrigEdge = 0
	EdgeRising TrigEdge = 1
)

func (p *Periph) SetTrigEdge(edge TrigEdge) {
	p.setTrigEdge(edge)
}

type Event uint16

const EvAll Event = evAll

type Error uint16

func (e Error) Error() string {
	return e.error()
}

const ErrAll Error = errAll

func (p *Periph) Status() (Event, Error) {
	return p.status()
}

func (p *Periph) Clear(ev Event, err Error) {
	p.clear(ev, err)
}

func (p *Periph) EnableIRQ(ev Event, err Error) {
	p.enableIRQ(ev, err)
}

func (p *Periph) DisableIRQ(ev Event, err Error) {
	p.disableIRQ(ev, err)
}

func (p *Periph) EnableDMA(circural bool) {
	p.enableDMA(circural)
}

func (p *Periph) DisableDMA() {
	p.disableDMA()
}

func (p *Periph) Start() {
	p.start()
}
