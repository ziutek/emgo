package main

import (
	"sync/atomic"

	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/ppi"
	"nrf5/hal/timer"
)

const (
	tachCC   = 0 // CC0 stores tach period.
	ovrfCC   = 1 // CC1 produces overflow event.
	ovrfStop = timer.COMPARE1_STOP
	ipr      = 3 //  Impulses per revolution.
	presc    = 4 // 1 MHz, overflow in 65.5 ms: min. freq 15.3 Hz (305 RPM).
)

type Tachometer struct {
	t        *timer.Periph
	te       gpiote.Chan
	ppiClear ppi.Channels
	ppiAll   ppi.Channels
	pins     [2]gpio.Pin // Support up to 2 sources. Can be replaced with slice.
	rpm      [2]uint32   // Support up to 2 sources. Can be replaced with slice.
	n, i     byte
}

func NewTachometer(t *timer.Periph, te gpiote.Chan, pc0, pc1, pc2, pc3 ppi.Chan, pg0, pg1 ppi.Group, pins ...gpio.Pin) *Tachometer {
	t.Task(timer.STOP).Trigger()
	t.StoreMODE(timer.TIMER)
	t.StoreBITMODE(timer.Bit16)
	t.StorePRESCALER(presc)
	t.Task(timer.CLEAR).Trigger()
	t.StoreCC(ovrfCC, 0)
	ev := t.Event(timer.COMPARE(ovrfCC))
	ev.Clear()
	ev.EnableIRQ()
	t.StoreSHORTS(ovrfStop)
	te.Setup(gpio.Pin{}, gpiote.ModeDiscon)

	// This code uses behavior undocumented for nRF51 but documented for nRF52
	// that EN task has priority over DIS task.
	ev = te.IN().Event()
	pc0.SetEEP(ev)
	pc0.SetTEP(t.Task(timer.CLEAR))
	pc1.SetEEP(ev)
	pc1.SetTEP(pg0.EN().Task())
	pc2.SetEEP(ev)
	pc2.SetTEP(t.Task(timer.CAPTURE(tachCC)))
	pc3.SetEEP(ev)
	pc3.SetTEP(pg1.DIS().Task())
	pg0.SetChannels(pc2.Mask() | pc3.Mask())
	pg1.SetChannels(pc0.Mask() | pc1.Mask() | pc2.Mask() | pc3.Mask())

	tach := new(Tachometer)
	tach.t = t
	tach.te = te
	tach.ppiClear = pc0.Mask() | pc1.Mask()
	tach.ppiAll = pc0.Mask() | pc1.Mask() | pc2.Mask() | pc3.Mask()
	for i, pin := range pins {
		pin.Setup(gpio.ModeIn)
		tach.pins[i] = pin
	}
	tach.n = byte(len(pins))
	tach.startMeasure()
	return tach
}

func (tach *Tachometer) startMeasure() {
	tach.ppiAll.Disable()
	tach.te.Setup(tach.pins[tach.i], gpiote.ModeEvent|gpiote.PolarityHiToLo)
	tach.t.StoreCC(tachCC, 0)
	tach.t.Task(timer.START).Trigger()
	tach.ppiClear.Enable()
}

// ISR returns the number of channel that has just been measured.
func (tach *Tachometer) ISR() int {
	tach.t.Event(timer.COMPARE(ovrfCC)).Clear()
	i := int(tach.i)
	cc := tach.t.LoadCC(tachCC)
	if cc != 0 {
		cc = 60 * 16e6 / (1 << presc * ipr) / cc
	}
	atomic.StoreUint32(&tach.rpm[i], tach.rpm[i]<<16|cc)
	if tach.i = byte(i + 1); tach.i == tach.n {
		tach.i = 0
	}
	tach.startMeasure()
	return i
}

func (tach *Tachometer) RPM(n int) int {
	rpm := atomic.LoadUint32(&tach.rpm[n])
	return int(rpm>>16+rpm&0xFFFF) / 2 // Avg. from previous and current RPM.
}
