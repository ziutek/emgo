package main

import (
	"delay"
	"sync/atomic"
)

var (
	waterSig = make(chan struct{}, 1)
	waterCnt int32
)

func waterIRQ() {
	atomic.AddInt32(&waterCnt, 1)
	select {
	case waterSig <- struct{}{}:
	default:
	}
	if waterCnt&1 != 0 {
		blue.Off()
	} else {
		blue.On()
	}
}

func waterTask() {
	for {
		if atomic.LoadInt32(&waterCnt) == 0 {
			blue.Off()
			<-waterSig
		}
		wf := int(atomic.SwapInt32(&waterCnt, 0)) * 16
		green.On()
		if wf > 500 {
			wf = 500
			beep(100)
			delay.Millisec(wf - 100)
		} else {
			delay.Millisec(wf)
		}
		green.Off()
		delay.Millisec(500 - wf)

	}
}
