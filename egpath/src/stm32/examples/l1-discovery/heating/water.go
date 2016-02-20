package main

import (
	"rtos"
	"sync/atomic"
)

var WaterHeaterOn int32

var water struct {
	Flag    rtos.EventFlag
	Counter int32
}

// waterISR is called for every pulse from water flow sensor.
func waterISR() {
	atomic.AddInt32(&water.Counter, 1)
	water.Flag.Set()
}

/*
func waterTask() {
	const (
		period = 500 // ms
		scale  = 40
	)
	r0, r1, r2 := W0, W1, W2
	for {
		if water.Flag.Val() == 0 {
			// Shuffle relays when idle to evenly use all heaters.
			r0, r1, r2 = r2, r0, r1
			atomic.StoreInt32(&WaterHeaterOn, 0)
			water.Flag.Wait(0)
		}
		water.Flag.Clear()
		ht := int(atomic.SwapInt32(&water.Counter, 0))
		if ht == 0 {
			continue
		}
		atomic.StoreInt32(&WaterHeaterOn, 1)

		start := rtos.Nanosec()
		ht *= scale // Heat time (if ht>period more than one heater is need).
		if ht > 3*period {
			// Only 3 heaters are connected.
			ht = 3 * period
			Blue.Set()
		}
		switch {
		case ht <= period:
			SSR.SetPins(r0)
			rtos.SleepUntil(start + int64(ht)*1e6)
			SSR.ClearPins(r0)
		case ht <= 2*period:
			SSR.SetPins(r0 | r1)
			rtos.SleepUntil(start + int64(ht-period)*1e6)
			SSR.ClearPins(r1)
		default:
			SSR.SetPins(r0 | r1 | r2)
			rtos.SleepUntil(start + int64(ht-2*period)*1e6)
			SSR.ClearPins(r2)
		}
		Blue.Clear()
		rtos.SleepUntil(start + period*1e6)
		SSR.ClearPins(r0 | r1 | r2)
	}
}
*/
