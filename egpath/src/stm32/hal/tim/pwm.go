package tim

import (
	"stm32/hal/system"
)

type PWM struct {
	P *Periph
}

// SetFreq setups input clock frequency of underlying timer to produce PWM
// waveform with period periodus miscroseconds. Max is a value that corresponds
// to 100% duty-cycle.
func (pwm PWM) SetFreq(periodus, max int) {
	p := pwm.P
	pclk := p.Bus().Clock()
	if pclk < system.AHB.Clock() {
		pclk *= 2
	}
	m := 1e6 * uint64(max)
	div := (uint64(pclk)*uint64(periodus) + m/2) / m
	p.PSC.Store(PSC(div - 1))
	p.ARR.Store(ARR(max - 1))
}
