package main

import (
	"sync/atomic"

	"delay"
	"onewire"

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
	t.CCER.SetBits(tim.CC1P)
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

type waterHeaterControl struct {
	pwm           PulsePWM3
	cnt           counter
	tempResp      chan int
	scale         int
	desiredTemp16 int // °C/16
	lastPWM       int

	TempSensor onewire.Dev
}

func (w *waterHeaterControl) DesiredTemp() int {
	return atomic.LoadInt(&w.desiredTemp16) / 16
}

func (w *waterHeaterControl) SetDesiredTemp(temp int) {
	atomic.StoreInt(&w.desiredTemp16, temp*16)
}

func (w *waterHeaterControl) LastPower() int {
	pwmMax := w.pwm.Max()
	return 24 * atomic.LoadInt(&w.lastPWM) / pwmMax
}

func (w *waterHeaterControl) Init(timPWM, timCnt *tim.TIM_Periph, pclk uint) {
	setupPulsePWM(timPWM, pclk, 500, 9999)
	w.pwm.Init(timPWM)
	w.cnt.Init(timCnt)
	w.tempResp = make(chan int, 1)
	w.SetDesiredTemp(40) // °C
	w.scale = w.pwm.Max() / 1300
}

var water waterHeaterControl

func waterCntISR() {
	water.cnt.ClearIF()
	water.pwm.Pulse()
}

func waterPWMISR() {
	water.pwm.ClearIF()
	cnt := water.cnt.LoadAndReset()

	const coldWaterTemp16 = 15 * 16 // Typical input water temp. [°C/16]
	desiredTemp16 := atomic.LoadInt(&water.desiredTemp16)
	delta16 := desiredTemp16 - coldWaterTemp16

	if dev := water.TempSensor; dev.Type() != 0 {
		select {
		case owd.Cmd <- TempCmd{Dev: dev, Resp: water.tempResp}:
		default:
		}
		select {
		case temp16 := <-water.tempResp:
			if temp16 == InvalidTemp {
				break
			}
			delta16 += desiredTemp16 - temp16
		default:
			ledGreen.Set()
			delay.Loop(5e4)
			ledGreen.Clear()
		}
	}
	if delta16 < 0 {
		delta16 = 0
	} else if delta16 > 50*16 {
		delta16 = 50 * 16
	}
	pwm16 := delta16 * cnt * water.scale
	if pwm16 < 0 {
		pwm16 = 0
		ledGreen.Set()
		delay.Loop(5e4)
		ledGreen.Clear()
	}
	pwm := pwm16 / 16
	water.pwm.Set(pwm)
	atomic.StoreInt(&water.lastPWM, pwm)
}
