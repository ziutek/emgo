package main

import (
	"delay"
	//"fmt"
	"rtos"
	"sync/fence"

	"stm32/hal/raw/tim"
)

// setupPWM setups all channels of timer t as PWM output.
func setupPWM(t *tim.TIM_Periph, pclk uint, freqHz, max int) {
	div := uint(freqHz * max)
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
	SR    int

	sample []byte
	end    rtos.EventFlag
	k      int
	prev   int
}

func (a *Audio) ISR() {
	t := a.Timer
	s := a.sample
	t.SR.Store(0)
	if a.k == len(s)*2 {
		t.DIER.Store(0)
		a.end.Signal(1)
	} else {
		// Linear interpolation.
		v := int(s[a.k/2])
		if a.k&1 == 0 {
			v = (a.prev + v) / 2
		} else {
			a.prev = v
		}
		t.CCR1.Store(tim.CCR1_Bits(v))
		a.k++
	}
}

func (a *Audio) Play(sample []byte, dly int) {
	a.sample = sample
	a.prev = int(sample[0])
	a.k = 0
	a.end.Reset(0)
	fence.RW()
	a.Timer.DIER.Store(tim.UIE)
	a.end.Wait(1, 0)
	delay.Millisec(dly)
}

var audio Audio
