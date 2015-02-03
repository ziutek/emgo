package main

import (
	"delay"
	"rtos"
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
	period = 500 // ms
	scale  = 100
	max    = 1000
)

func waterTask() {
	r0, r1, r2 := ssr0, ssr1, ssr2
	for {
		select {
		case <-waterSig:
		default:
			// Shuffle relays when iddle ot evenly use all heaters.
			r0, r1, r2 = r2, r0, r1
			<-waterSig
		}
		wf := int(atomic.SwapInt32(&waterCnt, 0))
		if wf == 0 {
			continue
		}

		end := rtos.Uptime() + period*1e6
		wf = wf * scale / period
		if wf > max {
			wf = max
			blue.On()
		}

		switch {
		case wf <= max/3:
			ssrPort.SetBits(1 << r0)
			delay.Millisec(period * wf / (max / 3))
			ssrPort.ClearBit(r0)
		case wf <= max*2/3:
			ssrPort.SetBits(1<<r0 | 1<<r1)
			delay.Millisec(period * (wf - max/3) / (max / 3))
			ssrPort.ClearBit(r1)
		default:
			ssrPort.SetBits(1<<r0 | 1<<r1 | 1<<r2)
			delay.Millisec(period * (wf - max*2/3) / (max / 3))
			ssrPort.ClearBit(r2)
		}

		blue.Off()
		rtos.SleepUntil(end)
		ssrPort.ClearBits(1<<r0 | 1<<r1 | 1<<r2)
	}
}
