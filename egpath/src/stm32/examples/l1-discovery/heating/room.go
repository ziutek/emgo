package main

import (
	"delay"
	"sync/atomic"
	"time"

	"stm32/hal/raw/tim"
)

const rhMax = 10000

func setupPWM(t *tim.TIM_Periph, pclk uint, periodms, arr int) {
	t.PSC.U16.Store(uint16(int(pclk/1000)*periodms/(arr+1) - 1))
	t.ARR.U32.Store(uint32(arr))
	t.CCMR1.Store(6<<tim.OC1Mn | 6<<tim.OC2Mn | tim.OC1PE | tim.OC2PE)
	t.CCMR2.Store(7<<tim.OC3Mn | 6<<tim.OC4Mn | tim.OC3PE | tim.OC4PE)
	t.CCR2.Store(0)
	t.CCR3.U32.Store(uint32(arr + 1))
	t.CCR4.Store(0)
	t.CCER.Store(tim.CC2E | tim.CC3E | tim.CC4E)
	t.CR1.Store(tim.URS | tim.CEN)
}

type roomHeaterControl struct {
	desiredTemp16 int // °C/16
	tempSensor    Sensor
}

func (r *roomHeaterControl) TempSensor() *Sensor {
	return &r.tempSensor
}

func (r *roomHeaterControl) DesiredTemp16() int {
	return atomic.LoadInt(&r.desiredTemp16)
}

func (r *roomHeaterControl) SetDesiredTemp16(temp16 int) {
	atomic.StoreInt(&r.desiredTemp16, temp16)
}

func (r *roomHeaterControl) loop(t *tim.TIM_Periph) {
	og := &t.CCR2.U32
	la := &t.CCR3.U32
	st := &t.CCR4.U32
	tempResp := make(chan int, 1)
	for {
		p := 0
		t := time.Now()
		hm := t.Hour()*60 + t.Minute()
		dev := r.tempSensor.Load()
		// Heat only if tempSensor set and cheap energy time: 22-6 and
		// 13-15: 1/10..29/02 or 15-17: 3/01..30/09.
		const (
			offset     = -37 // my electric meter is 37 minutes late.
			margin     = 5
			nightStart = 22*60 + margin + offset
			nightEnd   = 6*60 - margin + offset
		)
		dayStart := 13*60 + margin + offset
		dayEnd := 15*60 - margin + offset
		if m := t.Month(); m >= 3 && m < 10 {
			dayStart = 15*60 + margin + offset
			dayEnd = 17*60 - margin + offset
		}
		if t.Year() >= 2020 &&
			dev.Type() != 0 && (hm < nightEnd || hm > nightStart ||
			hm > dayStart && hm < dayEnd) {

			owd.Cmd <- TempCmd{Dev: dev, Resp: tempResp}
			t := <-tempResp
			if t != InvalidTemp {
				desiredTemp16 := atomic.LoadInt(&r.desiredTemp16)
				p = (desiredTemp16 - t) * rhMax / (2 * 16) // 1°C : maxP/2
				// Disable room heater if water heater exceeds 9 kW.
				maxP := (9 - water.LastPower()) * rhMax
				switch {
				case maxP < 0:
					maxP = 0
				case maxP > rhMax:
					maxP = rhMax
				}
				switch {
				case p < 0:
					p = 0
				case p > maxP:
					p = maxP
				}
			}
		}
		og.Store(uint32(p) * 2 / 3) // Small room, PWM mode 6.
		la.Store(rhMax - uint32(p)) // PWM mode 7 (use tail of PWM period first).
		st.Store(uint32(p) * 2 / 3) // Medium room, PWM mode 6.
		delay.Millisec(5e3)
	}
}

func (r *roomHeaterControl) Start(t *tim.TIM_Periph, pclk uint) {
	setupPWM(t, pclk, 500, rhMax-1)
	r.SetDesiredTemp16(21 * 16) // °C/16
	go r.loop(t)
}

var room roomHeaterControl
