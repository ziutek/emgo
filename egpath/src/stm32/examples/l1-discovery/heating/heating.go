package main

import (
	"delay"
	"rtos"
	"sync/atomic"
)

var hlevel, waterPrio int32

func buttonIRQ() {
	atomic.StoreInt32(&hlevel, (hlevel+1)&3)
}

func heatingTask() {
	for {
		rtos.Debug(0).WriteString("Hello debugger!\n")
		green.Blink(10)

		hl := atomic.LoadInt32(&hlevel)
		if hl == 0 || atomic.LoadInt32(&waterPrio) != 0 {
			heatPort.ClearBits(1<<heat0 | 1<<heat1 | 1<<heat2)
			delay.Millisec(500)
		} else if hl == 1 {
			heatPort.ClearAndSet(1<<(16+heat0) | 1<<heat1)
			delay.Millisec(167)
			heatPort.ClearAndSet(1<<(16+heat1) | 1<<heat2)
			delay.Millisec(167)
			heatPort.ClearAndSet(1<<(16+heat2) | 1<<heat0)
			delay.Millisec(167)
		} else if hl == 2 {
			heatPort.ClearAndSet(1<<(16+heat0) | 1<<heat1 | 1<<heat2)
			delay.Millisec(167)
			heatPort.ClearAndSet(1<<(16+heat1) | 1<<heat2 | 1<<heat0)
			delay.Millisec(167)
			heatPort.ClearAndSet(1<<(16+heat2) | 1<<heat0 | 1<<heat1)
			delay.Millisec(167)
		} else {
			heatPort.SetBits(1<<heat0 | 1<<heat1 | 1<<heat2)
			delay.Millisec(500)
		}
	}
}
