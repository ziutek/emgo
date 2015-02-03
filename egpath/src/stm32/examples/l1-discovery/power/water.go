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
}

const (
	period = 500
	max    = period * 3 / 4
)

func waterTask() {
	for {
		<-waterSig
		wf := int(atomic.SwapInt32(&waterCnt, 0))
		if wf == 0 {
			continue
		}
		wf *= 13
		if wf > max {
			wf = max
			blue.On()
		}
		ssrPort.SetBits(1<<ssr0 | 1<<ssr1 | 1<<ssr2)
		delay.Millisec(wf)
		ssrPort.ClearBits(1<<ssr0 | 1<<ssr1 | 1<<ssr2)
		blue.Off()
		delay.Millisec(period - wf)

	}
}
