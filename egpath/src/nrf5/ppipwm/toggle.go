// Package ppipwm provides implementations of Pulse Width Modulation that uses
// PPI to connect compare events of TIMER peripheral to GPIOTE tasks.
package ppipwm

import (
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/ppi"
	"nrf5/hal/timer"
)

// Toggle provides PWM implementation based on gpiote OUT task configured in
// toggle mode. It is designed specifically for nRF51 chips because of
// limitations of their GPIOTE peripheral. Better solutions exists for nRF52.
type Toggle struct {
	pin gpio.Pin
	gc  gpiote.Chan
	t   *timer.Periph
	max int
}

// Make returns configured Toggle implementation of PPI based PWM. It uses
// specified pin, GPIOTE channel gc, timer t, and two PPI channels pc0, pc1.
func MakeToggle(pin gpio.Pin, gc gpiote.Chan, t *timer.Periph, pc0, pc1 ppi.Chan) Toggle {
	var pwm Toggle
	pwm.pin = pin
	pwm.t = t
	pwm.gc = gc
	pin.Clear()
	pin.Setup(gpio.ModeOut)
	t.Task(timer.STOP).Trigger()
	t.StoreSHORTS(timer.COMPARE1_CLEAR)
	pc0.SetEEP(t.Event(timer.COMPARE0))
	pc0.SetTEP(gc.OUT())
	pc0.Enable()
	pc1.SetEEP(t.Event(timer.COMPARE1))
	pc1.SetTEP(gc.OUT())
	pc1.Enable()
	return pwm
}

// NewToggle provides convenient way to create heap allocated Toggle struct. See
// MakeToggle for more information.
func NewToggle(pin gpio.Pin, gc gpiote.Chan, t *timer.Periph, pc0, pc1 ppi.Chan) *Toggle {
	pwm := new(Toggle)
	*pwm = MakeToggle(pin, gc, t, pc0, pc1)
	return pwm
}

// SetFreq sets prescaler to 2^pre and period to period microseconds.
func (pwm *Toggle) SetFreq(pre, period int) {
	if uint(pre) > 9 {
		panic("pwm: bad prescaler")
	}
	if period < 10 {
		panic("pwm: bad period")
	}
	t := pwm.t
	t.StorePRESCALER(pre)
	div := uint32(1) << uint(pre)
	max := 16*uint32(period)/div - 1
	if max > 0xFFFF {
		panic("pwm: bad pre and/or period for 16-bit timer")
	}
	t.StoreCC(1, max)
	pwm.max = int(max)
}

//  Max returns value that represents 100% duty cycle.
func (pwm *Toggle) Max() int {
	return pwm.max
}

// SetDutyCycle sets PWM duty cycle to dc. To obtain duty cycle in percent use
// the following formula: dc% = 100% * dc / pwm.Max().
func (pwm *Toggle) SetDutyCycle(dc int) {
	pin := pwm.pin
	gc := pwm.gc
	t := pwm.t
	if dc >= pwm.max {
		pin.Set()
		gc.Setup(pin, 0)
		return
	}
	pin.Clear()
	gc.Setup(pin, 0)
	if dc == 0 {
		return
	}
	t.Task(timer.STOP).Trigger()
	t.Task(timer.CLEAR).Trigger()
	t.StoreCC(0, uint32(dc))
	gc.Setup(pin, gpiote.ModeTask|gpiote.PolarityToggle|gpiote.OutInitHigh)
	t.Task(timer.START).Trigger()
}
