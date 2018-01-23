package main

import (
	"sync/atomic"

	"nrf5/ppipwm"
)

type FanControl struct {
	pwm    *ppipwm.Toggle
	tach   *Tachometer
	setrpm [2]int
}

func (fc *FanControl) SetRPM(n, rpm int) {
	atomic.StoreInt(&fc.setrpm[n], rpm)
}

func (fc *FanControl) RPM(n int) int {
	return fc.tach.RPM(n)
}

func (fc *FanControl) TachISR() {
	n := fc.tach.ISR()
	rpm := fc.tach.RPM(n)
	setrpm := fc.setrpm[n]
	e := setrpm - rpm
	_ = e
}
