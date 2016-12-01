package main

import (
	"delay"

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

func startRoomHeating(t *tim.TIM_Periph, pclk uint) {
	setupPWM(t, pclk, 500, rhMax-1)
	go roomHeatingLoop(t)
}

var desiredEnvTemp int = 22 * 16

func roomHeatingLoop(t *tim.TIM_Periph) {
	og := &t.CCR2.U32
	la := &t.CCR3.U32
	st := &t.CCR4.U32
	//tempResp := make(chan int, 1)
	for {
		p := 0
		dt := readRTC()
		hm := dt.Hour()*60 + dt.Minute()
		if (dt != DateTime{}) && water.LastPower() < 4 && 
			(hm < 5*60+45 || hm > 21*60+50 || hm > 12*60+50 && hm < 14*60+45) {
			/*
				owd.Cmd <- TempCmd{Dev: envTempSensor, Resp: tempResp}
				t := <-menu.tempResp
				if t == InvalidTemp {
					delay.Millisec(5e3)
					continue
				}
			*/
			t := 20 * 16
			p = (desiredEnvTemp - t) * 3333 / 16 // 1°C corresponds to maxP/3
			switch {
			case p < 0:
				p = 0
			case p > rhMax:
				p = rhMax
			}
		}
		og.Store(uint32(p) / 3)
		la.Store(rhMax - uint32(p))
		st.Store(uint32(p) / 2)
		delay.Millisec(5e3)
	}
}
