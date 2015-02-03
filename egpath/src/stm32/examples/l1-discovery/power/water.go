package main

import (
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
	scale  = 40
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
		ht := int(atomic.SwapInt32(&waterCnt, 0))
		if ht == 0 {
			continue
		}

		start := rtos.Uptime()
		ht *= scale // Heat time (if ht>period more than one heater is need).
		if ht > 3*period {
			// Only 3 heaters are connected.
			ht = 3 * period
			blue.On()
		}

		switch {
		case ht <= period:
			ssrPort.SetBits(1 << r0)
			rtos.SleepUntil(start + uint64(ht)*1e6)
			ssrPort.ClearBit(r0)
		case ht <= 2*period:
			ssrPort.SetBits(1<<r0 | 1<<r1)
			rtos.SleepUntil(start + uint64(ht-period)*1e6)
			ssrPort.ClearBit(r1)
		default:
			ssrPort.SetBits(1<<r0 | 1<<r1 | 1<<r2)
			rtos.SleepUntil(start + uint64(ht-2*period)*1e6)
			ssrPort.ClearBit(r2)
		}

		blue.Off()
		rtos.SleepUntil(start + period*1e6)
		ssrPort.ClearBits(1<<r0 | 1<<r1 | 1<<r2)
	}
}
