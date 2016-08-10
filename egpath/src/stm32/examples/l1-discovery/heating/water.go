package main

import (
	"stm32/hal/raw/tim"
)

type counter struct {
	t *tim.TIM_Periph
}

func (c *counter) Init(t *tim.TIM_Periph) {
	c.t = t
	t.CKD().Store(2 << tim.CKDn)
	// Connect CC1 to TI1, setup input filter.
	t.CCMR1.StoreBits(tim.CC1S|tim.IC1F, 1<<tim.CC1Sn|0xf<<tim.IC1Fn)
	// Set falling edge detection, enable CC1.
	t.CCER.SetBits(tim.CC1P | tim.CC1E)
	// Set external clock mode 1, clock from filtered TI1.
	t.SMCR.StoreBits(tim.SMS|tim.TS, 7<<tim.SMSn|5<<tim.TSn)
	// Use CC2 to generate an interrupt after first count.
	t.CCR2.Store(1)
	t.DIER.Store(tim.CC2IE)
	t.CEN().Set()
}

func (c *counter) Load() int {
	return int(c.t.CNT.Load())
}

func (c *counter) LoadAndReset() int {
	cnt := int(c.t.CNT.Load())
	c.t.EGR.Store(tim.UG)
	return cnt
}

func (c *counter) ClearIF() {
	c.t.SR.Store(0)
}

var (
	waterPWM pulsePWM3
	waterCnt counter
	pulseNum int
)

func waterCntISR() {
	waterCnt.ClearIF()
	waterPWM.Pulse()
}

func waterPWMISR() {
	waterPWM.ClearIF()
	cnt := waterCnt.LoadAndReset()
	waterPWM.Set(cnt * 1000)
}
