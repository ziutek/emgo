// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package adc

import (
	"delay"
)

func (d *CircDriver) enable(calibrate bool) {
	p := d.p
	p.Enable()
	delay.Millisec(1) // TODO: Reduce this to Tstab (1 µs).
	if calibrate {
		p.Calibrate()
	}
}
