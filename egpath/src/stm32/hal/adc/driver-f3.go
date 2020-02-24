// +build f303xe

package adc

import (
	"delay"
)

func (d *Driver) enable(calibrate bool) {
	p := d.P
	if calibrate {
		p.Calibrate()
		if clkmode := p.ClockMode(); clkmode != 0 {
			delay.Loop(5 << (clkmode - 1))
		} else {
			delay.Millisec(1) // TODO: Be more accurate (shorter delay).
		}
	}
	d.watch(Ready, 0)
	p.Enable()
	d.done.Wait(1, 0)
}

