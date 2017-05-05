// Package ppipwm allows to produce PWM signal using PPI, TIMER and GPIOTE
// peripherals.
package ppipwm

import (
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/ppi"
	"nrf5/hal/timer"
)

// Toggle implementats three channel PWM usind GPIOTE OUT task configured in
// toggle mode. It is designed specifically for nRF51 chips because of
// limitations of their GPIOTE peripheral. Better solutions exists for nRF52.
//
// To produce desired PWM waveform Toggle configures a timer compare register,
// using simple algorithm, that usually generates glitches when used on working
// PWM channel. After that, PWM works without any CPU intervention and produces
// proper waveform until next duty cycle change. This is fine to drive LEDs but
// can cause trubles in case of receivers that rely on stable PWM frequency or
// phase (eg. some servos). Correct implementation is difficult and quite
// expensive in case of nRF51.
//
// Toggle cannot be used concurently by multiple gorutines without proper
// synchronisation.
type Toggle struct {
	t  *timer.Periph
	gc [3]gpiote.Chan
}

// MakeToggle returns configured Toggle implementation of PPI based PWM using
// timer t.
func MakeToggle(t *timer.Periph) Toggle {
	t.Task(timer.STOP).Trigger()
	t.StoreSHORTS(timer.COMPARE3_CLEAR)
	return Toggle{t: t}
}

// NewToggle provides convenient way to create heap allocated Toggle struct. See
// MakeToggle for more information.
func NewToggle(t *timer.Periph) *Toggle {
	pwm := new(Toggle)
	*pwm = MakeToggle(t)
	return pwm
}

// SetFreq sets prescaler to 2^log2pre and period to periodus microseconds. It
// allows to configure a duty cycle with resolution = 16 * periodus / 2^log2pre.
// SetFreq returns (resolution-1), which is a value that should be passed to
// SetDC/SetInvDC to produce PWM with 100% duty cycle.
func (pwm *Toggle) SetFreq(log2pre, periodus int) int {
	if uint(log2pre) > 9 {
		panic("pwm: bad prescaler")
	}
	if periodus < 1 {
		panic("pwm: bad period")
	}
	t := pwm.t
	t.StorePRESCALER(log2pre)
	div := uint32(1) << uint(log2pre)
	max := 16*uint32(periodus)/div - 1
	if max > 0xFFFF {
		panic("pwm: max>0xFFFF")
	}
	t.StoreCC(3, max)
	return int(max)
}

// Max returns a value that corresponds to 100% PWM duty cycle. See SetFreq for
// more information.
func (pwm *Toggle) Max() int {
	return int(pwm.t.LoadCC(3))
}

func checkOutput(n int) {
	if uint(n) > 2 {
		panic("pwm: bad output")
	}
}

// Setup setups n-th of three PWM channels. Each PWM channel uses one GPIOTE
// channel and two PPI channels.
func (pwm *Toggle) Setup(n int, pin gpio.Pin, gc gpiote.Chan, pc0, pc1 ppi.Chan) {
	checkOutput(n)
	pin.Clear()
	pin.Setup(gpio.ModeOut)
	gc.Setup(pin, 0)
	pwm.gc[n] = gc
	t := pwm.t
	pc0.SetEEP(t.Event(timer.COMPARE(n)))
	pc0.SetTEP(gc.OUT().Task())
	pc0.Enable()
	pc1.SetEEP(t.Event(timer.COMPARE(3)))
	pc1.SetTEP(gc.OUT().Task())
	pc1.Enable()
}

// SetDC sets a duty cycle for n-th PWM channel.
func (pwm *Toggle) SetDC(n, dc int) {
	checkOutput(n)
	gc := pwm.gc[n]
	t := pwm.t
	pin, _ := gc.Config()
	t.Task(timer.STOP).Trigger()
	switch {
	case dc <= 0:
		pin.Clear()
		gc.Setup(pin, 0)
		return
	case dc >= pwm.Max():
		pin.Set()
		gc.Setup(pin, 0)
		return
	}
	cfg := gpiote.ModeTask | gpiote.PolarityToggle
	t.Task(timer.CAPTURE(n)).Trigger()
	cnt := int(t.LoadCC(n))
	if cnt < dc {
		gc.Setup(pin, cfg|gpiote.OutInitHigh)
	} else {
		gc.Setup(pin, cfg|gpiote.OutInitLow)
	}
	t.StoreCC(n, uint32(dc))
	t.Task(timer.START).Trigger()
}

// SetInvDC sets a duty cycle for n-th PWM channel. It produces inverted
// waveform.
func (pwm *Toggle) SetInvDC(n, dc int) {
	checkOutput(n)
	gc := pwm.gc[n]
	t := pwm.t
	pin, _ := gc.Config()
	t.Task(timer.STOP).Trigger()
	switch {
	case dc <= 0:
		pin.Set()
		gc.Setup(pin, 0)
		return
	case dc >= pwm.Max():
		pin.Clear()
		gc.Setup(pin, 0)
		return
	}
	cfg := gpiote.ModeTask | gpiote.PolarityToggle
	t.Task(timer.CAPTURE(n)).Trigger()
	cnt := int(t.LoadCC(n))
	if cnt < dc {
		gc.Setup(pin, cfg|gpiote.OutInitLow)
	} else {
		gc.Setup(pin, cfg|gpiote.OutInitHigh)
	}
	t.StoreCC(n, uint32(dc))
	t.Task(timer.START).Trigger()
}
