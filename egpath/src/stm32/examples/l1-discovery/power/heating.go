package main

import (
	"arch/cortexm/debug/itm"
	"delay"
)

func heatingTask() {
	//debug.SetDEMCR(debug.DEMCR() | debug.TrcEna)
	//itm.SetCtrl(itm.ITMEna)
	//itm.StimEnable(0)
	p := itm.StimPort(0)
	for {
		p.Store8('^')
		p.WriteString("Hello!\n")
		heatPort.ClearAndSet(1<<(16+heat0) | 1<<heat1)
		delay.Millisec(500)
		heatPort.ClearAndSet(1<<(16+heat1) | 1<<heat0)
		delay.Millisec(500)
	}
}
