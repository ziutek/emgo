// Package ppipwm provides implementations of Pulse Width Modulation that uses
// PPI to connect compare events of TIMER peripheral to GPIOTE tasks.
package ppipwm

import (
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/ppi"
	"nrf5/hal/timer"
)

// Toggle implementats three channel PWM based on gpiote OUT task configured in
// toggle mode. It is designed specifically for nRF51 chips because of
// limitations of their GPIOTE peripheral. Better solutions exists for nRF52.
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

// SetFreq sets prescaler to 2^log2pre and period to periodus microseconds.
// SetFreq returns a value that represents 100% duty cycle.
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

//  Max returns a value that represents 100% duty cycle.
func (pwm *Toggle) Max() int {
	return int(pwm.t.LoadCC(3))
}

const badOutput = "pwm: bad output"

// Setup setups one of three PWM channels.
func (pwm *Toggle) Setup(n int, pin gpio.Pin, gc gpiote.Chan, pc0, pc1 ppi.Chan) {
	if uint(n) > 2 {
		panic(badOutput)
	}
	pin.Clear()
	pin.Setup(gpio.ModeOut)
	gc.Setup(pin, 0)
	pwm.gc[n] = gc
	t := pwm.t
	pc0.SetEEP(t.COMPARE(n))
	pc0.SetTEP(gc.OUT())
	pc0.Enable()
	pc1.SetEEP(t.COMPARE(3))
	pc1.SetTEP(gc.OUT())
	pc1.Enable()
}

// SetDutyCycle for one of three PWM channels.
func (pwm *Toggle) SetDutyCycle(n, dc int) {
	// This is simple implementation that works well with LEDs but sometimes
	// produces glitches that are detected by receivers that rely on stable PWM
	// frequency/phase (eg. servos).
	if uint(n) > 2 {
		panic(badOutput)
	}
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
	t.CAPTURE(n).Trigger()
	cnt := int(t.LoadCC(n))
	if cnt < dc {
		gc.Setup(pin, cfg|gpiote.OutInitHigh)
	} else {
		gc.Setup(pin, cfg|gpiote.OutInitLow)
	}
	t.StoreCC(n, uint32(dc))
	t.Task(timer.START).Trigger()
}
