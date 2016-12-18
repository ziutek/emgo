package main

import (
	"rtos"
	"sync/fence"

	"sound/samples/8b14k7/piano"

	"stm32/hal/raw/tim"
)

func setupAudioPWM(t *tim.TIM_Periph, pclk uint, sr, max int) {
	sr *= 2  // Perform 2 x oversampling.
	max *= 2 // Avoid rounding during linear interpolation.
	div := uint(sr * max)
	t.PSC.Store(tim.PSC_Bits((pclk+div/2)/div - 1))
	t.ARR.Store(tim.ARR_Bits(max - 1))
	t.CCMR1.Store(6<<tim.OC1Mn | tim.OC1PE)
	t.CCR1.Store(0)
	t.CCER.Store(tim.CC1E)
	t.DIER.Store(tim.UIE)
	t.CR1.Store(tim.ARPE | tim.URS | tim.CEN)
}

type Audio struct {
	Timer *tim.TIM_Periph

	snd  []byte
	end  rtos.EventFlag
	n    int
	prev int
}

func (a *Audio) ISR() {
	t := a.Timer
	snd := a.snd
	t.SR.Store(0)
	if a.n>>1 == len(snd) {
		t.DIER.Store(0)
		a.end.Signal(1)
	} else {
		// Linear interpolation.
		var v int
		if a.n&1 == 0 {
			s := int(snd[a.n>>1])
			v = a.prev + s
			a.prev = s
		} else {
			v = 2 * a.prev
		}
		t.CCR1.Store(tim.CCR1_Bits(v))
		a.n++
	}
}

func (a *Audio) Play(snd []byte) {
	a.snd = snd
	a.prev = int(snd[0])
	a.n = 0
	a.end.Reset(0)
	fence.RW()
	a.Timer.DIER.Store(tim.UIE)
	a.end.Wait(1, 0)
}

var audio Audio

var (
	c2  = (*[8192]byte)(&piano.C2)
	c2s = (*[8192]byte)(&piano.C2s)
	d2  = (*[8192]byte)(&piano.D2)
	d2s = (*[8192]byte)(&piano.D2s)
	e2  = (*[8192]byte)(&piano.E2)
	f2  = (*[8192]byte)(&piano.F2)
	f2s = (*[8192]byte)(&piano.F2s)
	g2  = (*[8192]byte)(&piano.G2)
	g2s = (*[8192]byte)(&piano.G2s)
	a2  = (*[8192]byte)(&piano.A2)
	a2s = (*[8192]byte)(&piano.A2s)
	h2  = (*[8192]byte)(&piano.H2)
	c3  = (*[8192]byte)(&piano.C3)
	c3s = (*[8192]byte)(&piano.C3s)
	d3  = (*[8192]byte)(&piano.D3)
	d3s = (*[8192]byte)(&piano.D3s)
	e3  = (*[8192]byte)(&piano.E3)
	f3  = (*[8192]byte)(&piano.F3)
	f3s = (*[8192]byte)(&piano.F3s)
	g3  = (*[8192]byte)(&piano.G3)
	g3s = (*[8192]byte)(&piano.G3s)
	a3  = (*[8192]byte)(&piano.A3)
	a3s = (*[8192]byte)(&piano.A3s)
	h3  = (*[8192]byte)(&piano.H3)
	c4  = (*[8192]byte)(&piano.C4)
)

type Note struct {
	Sample *[8192]byte
	Delay  int
}

var melody = [...]Note{
	{c2, 100},
	{c2s, 100},
	{d2, 100},
	{d2s, 100},
	{e2, 100},
	{f2, 100},
	{f2s, 100},
	{g2, 100},
	{g2s, 100},
	{a2, 100},
	{a2s, 100},
	{h2, 100},
	{c3, 100},
	{c3s, 400},

	{h2, 100},
	{g2, 100},
	{a2, 100},
	{a2, 100},

	{g2, 100},
	{g2, 100},
	{a2, 400},

	{g2, 100},
	{g2, 100},
	{a2, 100},
	{a2, 100},

	{h2, 100},
	{g2, 100},
	{a2, 400},

	{g2, 100},
	{g2, 100},
	{a2, 100},
	{a2, 100},

	{h2, 100},
	{g2, 100},
	{a2, 100},
	{a2, 100},

	{g2, 100},
	{g2, 100},
	{a2, 100},
	{a2, 100},

	{h2, 100},
	{g2, 100},
	{a2, 100},
	{a2, 100},

	{g2, 100},
	{g2, 100},
	{a2, 100},
	{a2, 100},
	{g2, 800},

	{g2, 100},
	{g2, 100},
	{a2, 100},
	{a2, 100},
	{g2, 800},
}
