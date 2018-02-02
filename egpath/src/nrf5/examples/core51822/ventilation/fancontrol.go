package main

import (
	"delay"
	"fmt"
	"sync/atomic"

	"nrf5/ppipwm"
)

const (
	pwmNum  = 35
	stepRPM = 1 << 5
	minRPM  = 420
	maxRPM  = minRPM + (pwmNum-1)*stepRPM
	divP    = 8
	divI    = 64
	divD    = 16
)

type fan struct {
	rpmToPWM  [pwmNum]byte
	targetRPM int
	sum       int
	lastE     int
}

func (f *fan) TargetRPM() int {
	return atomic.LoadInt(&f.targetRPM)
}

func (f *fan) SetTargetRPM(rpm int) {
	atomic.StoreInt(&f.targetRPM, rpm)
}

func (f *fan) NextE(e, maxSum int) (sum, diff int) {
	sum = f.sum + e
	if sum > maxSum {
		sum = maxSum
	} else if sum < -maxSum {
		sum = -maxSum
	}
	f.sum = sum
	diff = e - f.lastE
	f.lastE = e
	return sum, diff
}

func (f *fan) ResetE() {
	f.sum = 0
	f.lastE = 0
}

func (f *fan) ModelPWM(rpm int) int {
	r := rpm - minRPM
	n := r / stepRPM
	m := r & (stepRPM - 1)
	if n < 0 {
		return 0
	}
	if n >= len(f.rpmToPWM)-1 {
		return int(f.rpmToPWM[len(f.rpmToPWM)-1])
	}
	a := int(f.rpmToPWM[n])
	b := int(f.rpmToPWM[n+1])
	return ((stepRPM-m)*a + m*b) / stepRPM
}

func (f *fan) SetModelPWM(n, pwm int) {
	f.rpmToPWM[n] = byte(pwm)
}

type FanControl struct {
	pwm  *ppipwm.Toggle
	tach *Tachometer
	fans [2]fan
	maxI int
}

func NewFanControl(pwm *ppipwm.Toggle, tach *Tachometer) *FanControl {
	fc := new(FanControl)
	fc.pwm = pwm
	fc.tach = tach
	return fc
}

func (fc *FanControl) MaxRPM() int {
	return maxRPM
}

func (fc *FanControl) TargetRPM(n int) int {
	return fc.fans[n].TargetRPM()
}

func (fc *FanControl) SetTargetRPM(n, rpm int) {
	if rpm < 0 {
		rpm = 0
	} else if rpm > maxRPM {
		rpm = maxRPM
	}
	fc.fans[n].SetTargetRPM(rpm)
}

func (fc *FanControl) RPM(n int) int {
	return fc.tach.RPM(n)
}

func (fc *FanControl) TachISR() {
	n := fc.tach.ISR()
	fan := &fc.fans[n]
	targetRPM := fan.TargetRPM()
	if targetRPM < 0 {
		return
	}
	dc := 0
	if targetRPM >= minRPM {
		modelPWM := fan.ModelPWM(targetRPM)
		rpm := fc.RPM(n)
		e := targetRPM - rpm
		sum, diff := fan.NextE(e, fc.maxI)
		dc = modelPWM + e/divP + sum/divI + diff/divD
	} else {
		fan.ResetE()
	}
	fc.pwm.SetInv(n, dc)
}

func (fc *FanControl) Identify() {
	for n := range fc.fans {
		fan := &fc.fans[n]
		fan.SetTargetRPM(-1) // Prevent useing PWM by TachISR.
		fc.pwm.SetInv(n, 0)
		for i := range fan.rpmToPWM {
			fan.rpmToPWM[i] = 0
		}
	}
	maxPWM := fc.pwm.Max()
	if maxPWM > 255 {
		panic("maxPWM>255")
	}
	fc.maxI = maxPWM * divI / 2
	todo := uint(1<<uint(len(fc.fans)) - 1)
	for pwm := 33; pwm < maxPWM && todo != 0; pwm++ {
		fc.pwm.SetManyInv(todo, pwm, pwm, pwm)
		delay.Millisec(500)
		for n := range fc.fans {
			fanMask := uint(1 << uint(n))
			if todo&fanMask == 0 {
				continue
			}
			rpm := fc.RPM(n)
			fmt.Printf("fan%d: pwm=%d rpm=%d\n", n, pwm, rpm)
			m := (rpm - minRPM + stepRPM - 1) / stepRPM
			switch {
			case m >= pwmNum:
				todo &^= 1 << uint(n)
				fc.pwm.SetInv(n, 0)
			case m >= 0:
				fc.fans[n].SetModelPWM(m, pwm)
			}
		}
	}
	fc.pwm.SetManyInv(todo, 0, 0, 0)
	for n := range fc.fans {
		fmt.Printf("fan%d:\n", n)
		fan := &fc.fans[n]
		fan.SetTargetRPM(0)
		for rpm := minRPM; rpm <= maxRPM; rpm += stepRPM / 2 {
			fmt.Printf("rpm=%d modelPWM=%d\n", rpm, fan.ModelPWM(rpm))
		}
	}
	// TODO: Use ModelPWM and todo to detect broken fan.
}
