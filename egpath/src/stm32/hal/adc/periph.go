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

type ClockMode byte

const (
	CKADC ClockMode = 0 // Asynchronous clock
	HCLK1 ClockMode = 1 // AHB clock / 1.
	HCLK2 ClockMode = 2 // AHB clock / 2.
	HCLK4 ClockMode = 3 // AHB clock / 4.
)

func (p *Periph) SetClockMode(ckmode ClockMode) {
	p.common().CKMODE().Store(adc.CCR_Bits(ckmode) << adc.CKMODEn)
}

// Calibrate calibrates p.
func (p *Periph) Callibrate() {
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

type SmplTime int8

// MaxSmplTime returns largest possible value of sampling time that takes no
// more than ht half clock cycles. If ht is too small MaxSmplTime returns
// negative value.
func MaxSmplTime(hc int) SmplTime {
	var i int
	for i < len(halfCycles) {
		if int(halfCycles[i]) > hc {
			break
		}
		i++
	}
	return SmplTime(i - 1)
}

func (st SmplTime) HalfCycles() int {
	return int(halfCycles[st])
}

func (p *Periph) SetSmplTime(ch int, st SmplTime) {
	p.setSmplTime(ch, st)
}

func (p *Periph) SetRegularSeq(ch ...int) {
	p.setRegularSeq(ch)
}

type TrigSrc byte

// SetTrigSrc selects source of external trigger.
func (p *Periph) SetTrigSrc(src TrigSrc) {
	p.setTrigSrc(src)
}

type TrigEdge byte

const (
	EdgeNone    TrigEdge = 0
	EdgeRising  TrigEdge = 1
	EdgeFalling TrigEdge = 2
)

func (p *Periph) SetTrigEdge(edge TrigEdge) {
	p.raw.EXTEN().Store(adc.CFGR_Bits(edge) << adc.EXTENn)
}

type Event uint16

const EvAll = evAll

func (p *Periph) Event() Event {
	return p.event()
}

func (p *Periph) Clear(events Event) {
	p.clear(events)
}

func (p *Periph) EnableIRQ(events Event) {
	p.enableIRQ(events)
}

func (p *Periph) DisableIRQ(events Event) {
	p.disableIRQ(events)
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

func panicCN() {
	panic("adc: bad channel number")
}
