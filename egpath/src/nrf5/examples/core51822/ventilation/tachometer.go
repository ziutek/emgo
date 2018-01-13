package main

import (
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/ppi"
	"nrf5/hal/timer"
)

type Tachometer timer.Periph

func MakeTachometer(t *timer.Periph, pin gpio.Pin, te gpiote.Chan, pp0, pp1 ppi.Chan) *Tachometer {
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
	return (*Tachometer)(t)
}

func (tach *Tachometer) RPM() int {
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
