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
	t.Task(timer.CLEAR).Trigger()
	t.StoreMODE(timer.TIMER)
	t.StoreBITMODE(timer.Bit16)
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

// SetFreq sets timer prescaler to 2^log2pre and period to periodus microseconds
// (log2pre must be in range [0..9]). It allows to configure a duty cycle with
// a resolution = 16 * periodus / 2^log2pre. Toggle uses timer in 16-bit mode so
// the resolution must be <= 65536. SetFreq returns (resolution-1), which is a
// value that should be passed to SetDuty/SetInvDuty to produce PWM with 100%
// duty cycle.
func (pwm *Toggle) SetFreq(log2pre, periodus int) int {
	if uint(log2pre) > 9 {
		panic("pwm: bad prescaler")
	}
	if periodus < 1 {
		panic("pwm: bad period")
	}
	t := pwm.t
	t.StorePRESCALER(log2pre)
	max := 16*uint32(periodus)>>uint(log2pre) - 1 // 16 MHz * period / pre - 1
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

func checkChan(n int) {
	if uint(n) > 2 {
		panic("pwm: bad channel")
	}
}

// Setup setups n-th of three PWM channels. Each PWM channel uses one GPIOTE
// channel and two PPI channels.
func (pwm *Toggle) Setup(n int, pin gpio.Pin, gc gpiote.Chan, pc0, pc1 ppi.Chan) {
	checkChan(n)
	pin.Clear()
	pin.Setup(gpio.ModeOut)
	gc.Setup(pin, gpiote.ModeDiscon)
	pwm.gc[n] = gc
	t := pwm.t
	pc0.SetEEP(t.Event(timer.COMPARE(n)))
	pc0.SetTEP(gc.OUT().Task())
	pc0.Enable()
	pc1.SetEEP(t.Event(timer.COMPARE(3)))
	pc1.SetTEP(gc.OUT().Task())
	pc1.Enable()
}

// DutyCycle returns the current duty cycle on channel n. There is no way to
// check does it corresponds straight or inverted waveform.
func (pwm *Toggle) DutyCycle(n int) int {
	checkChan(n)
	return int(pwm.t.LoadCC(n))
}

// Set sets a duty cycle for n-th PWM channel. If dc > 0 or dc < pwm.Max()
// it stops the PWM timer and starts it just before return (this can produce
// glitches and affects all PWM channels).
func (pwm *Toggle) Set(n, dc int) {
	checkChan(n)
	gc := pwm.gc[n]
	t := pwm.t
	pin, _ := gc.Config()
	switch {
	case dc <= 0:
		pin.Clear()
		gc.Setup(pin, gpiote.ModeDiscon)
		return
	case dc >= pwm.Max():
		pin.Set()
		gc.Setup(pin, gpiote.ModeDiscon)
		return
	}
	t.Task(timer.STOP).Trigger()
	t.Task(timer.CAPTURE(n)).Trigger()
	cfg := gpiote.ModeTask | gpiote.PolarityToggle
	if int(t.LoadCC(n)) < dc {
		gc.Setup(pin, cfg|gpiote.OutInitHigh)
	} else {
		gc.Setup(pin, cfg|gpiote.OutInitLow)
	}
	t.StoreCC(n, uint32(dc))
	t.Task(timer.START).Trigger()
}

// SetInv works like Set but produces inverted waveform.
func (pwm *Toggle) SetInv(n, dc int) {
	checkChan(n)
	gc := pwm.gc[n]
	t := pwm.t
	pin, _ := gc.Config()
	switch {
	case dc <= 0:
		pin.Set()
		gc.Setup(pin, gpiote.ModeDiscon)
		return
	case dc >= pwm.Max():
		pin.Clear()
		gc.Setup(pin, gpiote.ModeDiscon)
		return
	}
	t.Task(timer.STOP).Trigger()
	t.Task(timer.CAPTURE(n)).Trigger()
	cfg := gpiote.ModeTask | gpiote.PolarityToggle
	if int(t.LoadCC(n)) < dc {
		gc.Setup(pin, cfg|gpiote.OutInitLow)
	} else {
		gc.Setup(pin, cfg|gpiote.OutInitHigh)
	}
	t.StoreCC(n, uint32(dc))
	t.Task(timer.START).Trigger()
}

// SetMany sets a duty cycles for PWM channels specifid by mask. Use it for more
// than one channel to minimizes the number of times the PWM timer is stopped
// and started (should produce less glitches than call SetDuty multiple times).
func (pwm *Toggle) SetMany(mask uint, dc0, dc1, dc2 int) {
	t := pwm.t
	t.Task(timer.STOP).Trigger()
	for n, dc := range [3]int{dc0, dc1, dc2} {
		if mask&1 != 0 {
			gc := pwm.gc[n]
			pin, _ := gc.Config()
			switch {
			case dc <= 0:
				pin.Clear()
				gc.Setup(pin, gpiote.ModeDiscon)
			case dc >= pwm.Max():
				pin.Set()
				gc.Setup(pin, gpiote.ModeDiscon)
			default:
				t.Task(timer.CAPTURE(n)).Trigger()
				cfg := gpiote.ModeTask | gpiote.PolarityToggle
				if int(t.LoadCC(n)) < dc {
					gc.Setup(pin, cfg|gpiote.OutInitHigh)
				} else {
					gc.Setup(pin, cfg|gpiote.OutInitLow)
				}
				t.StoreCC(n, uint32(dc))
			}
		}
		mask >>= 1
	}
	t.Task(timer.START).Trigger()
}

// SetManyInv works like SetMany but produces inverted waveform.
func (pwm *Toggle) SetManyInv(mask uint, dc0, dc1, dc2 int) {
	t := pwm.t
	t.Task(timer.STOP).Trigger()
	for n, dc := range [3]int{dc0, dc1, dc2} {
		if mask&1 != 0 {
			gc := pwm.gc[n]
			pin, _ := gc.Config()
			switch {
			case dc <= 0:
				pin.Set()
				gc.Setup(pin, gpiote.ModeDiscon)
			case dc >= pwm.Max():
				pin.Clear()
				gc.Setup(pin, gpiote.ModeDiscon)
			default:
				t.Task(timer.CAPTURE(n)).Trigger()
				cfg := gpiote.ModeTask | gpiote.PolarityToggle
				if int(t.LoadCC(n)) < dc {
					gc.Setup(pin, cfg|gpiote.OutInitLow)
				} else {
					gc.Setup(pin, cfg|gpiote.OutInitHigh)
				}
				t.StoreCC(n, uint32(dc))
			}
		}
		mask >>= 1
	}
	t.Task(timer.START).Trigger()
}
