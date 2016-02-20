package main

import (
	"rtos"
	"sync/atomic"
)

var (
	waterSig = make(chan struct{}, 1)
	waterCnt int32
)

// waterIRQ is called for every pulse from water flow sensor.
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
			// Shuffle relays when iddle to evenly use all heaters.
			r0, r1, r2 = r2, r0, r1
			atomic.StoreInt32(&waterPrio, 0)
			<-waterSig
		}
		ht := int(atomic.SwapInt32(&waterCnt, 0))
		if ht == 0 {
			continue
		}
		atomic.StoreInt32(&waterPrio, 1)

		start := rtos.Uptime()
		ht *= scale // Heat time (if ht>period more than one heater is need).
		if ht > 3*period {
			// Only 3 heaters are connected.
			ht = 3 * period
			blue.On()
		}

		switch {
		case ht <= period:
			ssrPort.SetPin(r0)
			rtos.SleepUntil(start + int64(ht)*1e6)
			ssrPort.ClearPin(r0)
		case ht <= 2*period:
			ssrPort.SetPins(1<<uint(r0) | 1<<uint(r1))
			rtos.SleepUntil(start + int64(ht-period)*1e6)
			ssrPort.ClearPin(r1)
		default:
			ssrPort.SetPins(1<<uint(r0) | 1<<uint(r1) | 1<<uint(r2))
			rtos.SleepUntil(start + int64(ht-2*period)*1e6)
			ssrPort.ClearPin(r2)
		}

		blue.Off()
		rtos.SleepUntil(start + period*1e6)
		ssrPort.ClearPins(1<<uint(r0) | 1<<uint(r1) | 1<<uint(r2))
	}
}
