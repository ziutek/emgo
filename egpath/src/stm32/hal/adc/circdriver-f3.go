// +build f303xe

package adc

import (
	"delay"
)

func (d *CircDriver) enable(calibrate bool) {
	p := d.p
	if calibrate {
		p.Calibrate()
		if clkmode := p.ClockMode(); clkmode != 0 {
			delay.Loop(5 << (clkmode - 1))
		} else {
			delay.Millisec(1) // TODO: Be more accurate (shorter delay).
		}
	}
	d.watch(Ready)
	p.Enable()
	<-d.hc
}
