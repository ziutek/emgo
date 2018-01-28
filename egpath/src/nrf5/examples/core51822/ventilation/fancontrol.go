package main

import (
	"sync/atomic"

	"nrf5/ppipwm"
)

const pwmNum = 128

type fan struct {
	pwm [pwmNum]byte
	rpm int
}

func (f *fan) RPM() int {
	return f.rpm
}

func (f *fan) AtomicSetRPM(rpm int) {
	atomic.StoreInt(&f.rpm, rpm)
}

func (f *fan) PWM() int {
	r := f.rpm - 400
	n := r / 16
	m := r & 15
	if n < 0 {
		return 0
	}
	if n >= len(f.pwm)-1 {
		return int(f.pwm[len(f.pwm)-1])
	}
	a := int(f.pwm[n])
	b := int(f.pwm[n+1])
	return ((16-m)*a + m*b) / 16
}

func (f *fan) TestRPM(n int) int {
	if n >= pwmNum {
		return -1
	}
	return 400 + n*16
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
	return 400 + (pwmNum-1)*16
}

func (fc *FanControl) SetRPM(n, rpm int) {
	fc.fan[n].AtomicSetRPM(rpm)
}

func (fc *FanControl) RPM(n int) int {
	return fc.tach.RPM(n)
}

func (fc *FanControl) TachISR() {
	n := fc.tach.ISR()
	rpm := fc.tach.RPM(n)
	setrpm := fc.fan[n].RPM()
	e := setrpm - rpm
	_ = e
}
