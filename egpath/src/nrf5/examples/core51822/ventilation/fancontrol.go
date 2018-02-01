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

func (f *fan) TestRPM(n int) int {
	return minRPM + n*stepRPM
}

func (f *fan) SetModelPWM(n, pwm int) {
	f.rpmToPWM[n] = byte(pwm)
}

type FanControl struct {
	pwm  *ppipwm.Toggle
	tach *Tachometer
	fan  [2]fan
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
	return fc.fan[n].TargetRPM()
}

func (fc *FanControl) SetTargetRPM(n, rpm int) {
	if rpm < 0 {
		rpm = 0
	} else if rpm > maxRPM {
		rpm = maxRPM
	}
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
	for n := range fc.fan {
		fc.fan[n].SetTargetRPM(-1) // Prevent useing PWM by TachISR.
	}
	for n := range fc.fan {
		fc.pwm.SetInv(n, 0)
	}
	maxPWM := fc.pwm.Max()
	fc.maxI = maxPWM * divI / 2
	for n := 1; n < 2; n++ {
		fan := &fc.fan[n]
		pwm := 33
		lastRPM := 0
		for k := 0; k < pwmNum; k++ {
			testRPM := fan.TestRPM(k)
			for {
				fc.pwm.SetInv(n, pwm)
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
			fmt.Printf(" (%d) rpmToPWM[%d]=%d\n", testRPM, k, fan.rpmToPWM[k])
		}
		fc.pwm.SetInv(n, 0)
	}
	for n := range fc.fan {
		fc.fan[n].SetTargetRPM(0)
	}
}
