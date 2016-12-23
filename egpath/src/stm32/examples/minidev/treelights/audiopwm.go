package main

import (
	"rtos"
	"sync/fence"

	"stm32/hal/system"

	"stm32/hal/raw/tim"
)

type Audio struct {
	t    *tim.TIM_Periph
	snd  []byte
	end  rtos.EventFlag
	n    int
	zero int
}

// Setup setups audio PWM using t for sr symbol rate (Bd/s) and signed samples
// in range [-max, max].
func (a *Audio) Setup(t *tim.TIM_Periph, pclk uint, sr, max int) {
	if pclk != system.AHB.Clock() {
		pclk *= 2
	}
	a.t = t
	a.zero = max
	max *= 2 // [-max, max]
	sr *= 2  // Perform 2 x oversampling.
	div := uint(sr * max)
	t.PSC.Store(tim.PSC_Bits((pclk+div/2)/div - 1))
	t.ARR.Store(tim.ARR_Bits(max - 1)) // CCR=0: PWM=0%, CCR=max: PWM=100%.
	t.CCMR2.Store(6<<tim.OC3Mn | tim.OC3PE | 6<<tim.OC4Mn | tim.OC4PE)
	t.CCR3.Store(0)
	t.CCR4.Store(0)
	t.CCER.Store(tim.CC3E | tim.CC4E)
	t.DIER.Store(tim.UIE)
	t.CR1.Store(tim.ARPE | tim.URS | tim.CEN)
}

func (a *Audio) ISR() {
	a.t.SR.Store(0)
	snd := a.snd
	n := a.n >> 1
	if n == len(snd) {
		a.t.DIER.Store(0)
		a.end.Signal(1)
		return
	}
	// Oversampling, 15-tap low-pass FIR filter: http://t-filter.engineerjs.com/
	var v int
	if a.n&1 == 0 {
		v = 1549 * (int(int8(snd[n])) + int(int8(snd[n-7])))
		v += 6936 * (int(int8(snd[n-1])) + int(int8(snd[n-6])))
		v += -7428 * (int(int8(snd[n-2])) + int(int8(snd[n-5])))
		v += 20929 * (int(int8(snd[n-3])) + int(int8(snd[n-4])))
	} else {
		v = 6227 * (int(int8(snd[n])) + int(int8(snd[n-6])))
		v += -299 * (int(int8(snd[n-1])) + int(int8(snd[n-5])))
		v += 386 * (int(int8(snd[n-2])) + int(int8(snd[n-4])))
		v += 32338 * int(int8(snd[n-3]))
	}
	v >>= 15
	// Differential output.
	a.t.CCR3.Store(tim.CCR3_Bits(a.zero + v))
	a.t.CCR4.Store(tim.CCR4_Bits(a.zero - v))
	a.n++
}

// Play plays snd. Snd should contain signed (int8) values.
func (a *Audio) Play(snd []byte) {
	a.snd = snd
	a.n = 7 << 1 // Skip fierst 7 samples to speed up FIR algorithm.
	a.end.Reset(0)
	fence.RW()
	a.t.DIER.Store(tim.UIE)
	a.end.Wait(1, 0)
}

var audio Audio
