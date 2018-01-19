package main

import (
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/ppi"
	"nrf5/hal/rtc"
	"nrf5/hal/timer"
)

// Tachometer supports up to 4 tach sources.
type Tachometer struct {
	cnt    *timer.Periph
	rtc    *rtc.Periph
	pins   [4]gpio.Pin
	te     gpiote.Chan
	ppi    ppi.Chan
	period uint32
	delay  uint16
	ccn    byte
	k, n   byte
}

func NewTachometer(cnt *timer.Periph, te gpiote.Chan, ppi ppi.Chan, rt *rtc.Periph, ccn, periodms int, pins ...gpio.Pin) *Tachometer {
	for _, pin := range pins {
		pin.Setup(gpio.ModeIn)
	}
	cnt.Task(timer.STOP).Trigger()
	cnt.StoreMODE(timer.COUNTER)
	cnt.StoreBITMODE(timer.Bit16)
	cnt.Task(timer.START).Trigger()
	ppi.SetEEP(te.IN().Event())
	ppi.SetTEP(cnt.Task(timer.COUNT))
	ppi.Enable()
	t := new(Tachometer)
	t.cnt = cnt
	for i, pin := range pins {
		t.pins[i] = pin
	}
	t.n = byte(len(pins))
	t.te = te
	t.ppi = ppi
	t.rtc = rt
	t.ccn = byte(ccn)
	presc := rt.LoadPRESCALER() + 1
	delay := 32768 * uint32(periodms) / (presc * 1e3 * uint32(len(pins)))
	t.delay = uint16(delay)
	t.period = delay * presc
	rt.StoreCC(ccn, rt.LoadCOUNTER()) // Avoid spurious interrupt.
	ev := rt.Event(rtc.COMPARE(ccn))
	ev.Clear()
	t.setupChan(0)
	ev.EnableIRQ()
	return t
}

func (t *Tachometer) setupChan(k int) {
	t.te.Setup(t.pins[k], gpiote.ModeEvent|gpiote.PolarityToggle)
	t.cnt.Task(timer.CLEAR).Trigger()
	t.rtc.StoreCC(int(t.ccn), t.rtc.LoadCOUNTER()+uint32(t.delay))
}

func (t *Tachometer) RTCISR() bool {
	ev := t.rtc.Event(rtc.COMPARE(int(t.ccn)))
	if !ev.IsSet() {
		return false
	}
	ev.Clear()
	cnt, k := t.cnt, int(t.k)
	cnt.Task(timer.CAPTURE(k)).Trigger()
	if k++; k == int(t.n) {
		k = 0
	}
	t.setupChan(k)
	t.k = byte(k)
	return k == 0
}

func (t *Tachometer) RPM(n int) int {
	const ipr = 6 // Impulses per revolution (both edges).
	return int(t.cnt.LoadCC(n) * 60 / ipr * 32768 / t.period)
}

// TachFast works well but uses too many resources.
type TachFast timer.Periph

func MakeTachFast(t *timer.Periph, pin gpio.Pin, te gpiote.Chan, pp0, pp1 ppi.Chan) *TachFast {
	pin.Setup(gpio.ModeIn)
	te.Setup(pin, gpiote.ModeEvent|gpiote.PolarityHiToLo)
	t.Task(timer.STOP).Trigger()
	t.StoreMODE(timer.TIMER)
	t.StoreBITMODE(timer.Bit16)
	t.StoreCC(0, 0)
	t.StoreCC(1, 0)
	t.StorePRESCALER(7) // 125 kHz
	pp0.SetEEP(te.IN().Event())
	pp0.SetTEP(t.Task(timer.CAPTURE(0)))
	pp0.Enable()
	pp1.SetEEP(te.IN().Event())
	pp1.SetTEP(t.Task(timer.CLEAR))
	pp1.Enable()
	t.Task(timer.START).Trigger()
	return (*TachFast)(t)
}

func (tach *TachFast) RPM() int {
	t := (*timer.Periph)(tach)
	cc := int(t.LoadCC(0))
	if ev := t.Event(timer.COMPARE(1)); ev.IsSet() {
		ev.Clear()
		t.StoreCC(0, 0)
		return 0
	}
	if cc == 0 {
		return 0
	}
	const ipr = 3 // Impulses per revolution.
	return 60 * 125e3 / ipr / cc
}
