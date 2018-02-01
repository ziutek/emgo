package main

import (
	"delay"
	"fmt"
	"sync/atomic"

	"nrf5/ppipwm"
)

const pwmNum = 32

type fan struct {
	rpmToPWM  [pwmNum]byte
	targetRPM int
}

func (f *fan) TargetRPM() int {
	return atomic.LoadInt(&f.targetRPM)
}

func (f *fan) SetTargetRPM(rpm int) {
	atomic.StoreInt(&f.targetRPM, rpm)
}

func (f *fan) ModelPWM(rpm int) int {
	r := rpm - 400
	n := r / 32
	m := r & 31
	if n < 0 {
		return 0
	}
	if n >= len(f.rpmToPWM)-1 {
		return int(f.rpmToPWM[len(f.rpmToPWM)-1])
	}
	a := int(f.rpmToPWM[n])
	b := int(f.rpmToPWM[n+1])
	return ((32-m)*a + m*b) / 32
}

func (f *fan) TestRPM(n int) int {
	return 400 + n*32
}

func (f *fan) SetModelPWM(n, pwm int) {
	f.rpmToPWM[n] = byte(pwm)
}

type FanControl struct {
	pwm  *ppipwm.Toggle
	tach *Tachometer
	fan  [2]fan
}

func NewFanControl(pwm *ppipwm.Toggle, tach *Tachometer) *FanControl {
	fc := new(FanControl)
	fc.pwm = pwm
	fc.tach = tach
	return fc
}

func (fc *FanControl) MaxRPM() int {
	return 400 + (pwmNum-1)*32
}

func (fc *FanControl) TargetRPM(n int) int {
	return fc.fan[n].TargetRPM()
}

func (fc *FanControl) SetTargetRPM(n, rpm int) {
	fc.fan[n].SetTargetRPM(rpm)
}

func (fc *FanControl) RPM(n int) int {
	return fc.tach.RPM(n)
}

func (fc *FanControl) TachISR() {
	n := fc.tach.ISR()
	fan := &fc.fan[n]
	targetRPM := fan.TargetRPM()
	if targetRPM < 0 {
		return
	}
	modelPWM := fan.ModelPWM(targetRPM)
	rpm := fc.RPM(n)
	e := targetRPM - rpm
	fc.pwm.SetInvVal(n, modelPWM+e/4)
}

func (fc *FanControl) Identify() {
	// ppipwm.Toggle cannot be used concurently. Prevent useing by TachISR.
	for n := range fc.fan {
		fc.fan[n].SetTargetRPM(-1)
	}
	for n := range fc.fan {
		fc.pwm.SetInvVal(n, 0)
	}
	maxPWM := fc.pwm.Max()
	for n := 1; n < 2; n++ {
		fan := &fc.fan[n]
		pwm := 38
		lastRPM := 0
		for k := 0; k < pwmNum; k++ {
			testRPM := fan.TestRPM(k)
			for {
				fc.pwm.SetInvVal(n, pwm)
				delay.Millisec(500)
				rpm := fc.RPM(n)
				fmt.Printf(" %d: %d", pwm, rpm)
				if rpm >= testRPM {
					if rpm-testRPM < testRPM-lastRPM {
						fan.SetModelPWM(k, pwm)
					} else {
						fan.SetModelPWM(k, pwm-1)
					}
					pwm++
					lastRPM = rpm
					break
				}
				if pwm == maxPWM {
					for k++; k < pwmNum; k++ {
						fan.SetModelPWM(k, maxPWM)
					}
					break
				}
				pwm++
				lastRPM = rpm
				fmt.Printf("\n")
			}
			fmt.Printf(
				" (%d) rpmToPWM[%d]=%d\n",
				testRPM, k, fan.rpmToPWM[k],
			)
		}
		fc.pwm.SetInvVal(n, 0)
	}
	for n := range fc.fan {
		fc.fan[n].SetTargetRPM(0)
	}
}
