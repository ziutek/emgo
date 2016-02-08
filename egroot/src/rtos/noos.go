// +build noos

package rtos

import (
	"syscall"
)

func sleepUntil(end int64) {
	for Nanosec() < end {
		syscall.SetAlarm(end)
		syscall.Alarm.Wait()
	}
}

type event struct {
	raw syscall.Event
}

func newEvent() *Event {
	e := new(Event)
	e.raw = syscall.AssignEventFlag() | syscall.Alarm | 1
	return e
}

func eventSend(e *Event) {
	raw := syscall.AtomicLoadEvent(&e.raw) &^ 1
	syscall.AtomicStoreEvent(&e.raw, raw)
	raw.Send()
}

func eventWaitUntil(e *Event, t int64) bool {
	state := syscall.AtomicLoadEvent(&e.raw)
	state0, state1 := state&^1, state|1
	for {
		if syscall.AtomicCompareAndSwapEvent(&e.raw, state0, state1) {
			return true
		}
		if t >= 0 && Nanosec() >= t {
			return false
		}
		syscall.SetAlarm(t)
		state0.Wait()
	}
}

func waitEvent(t int64, events []*Event) uint32 {
	var sum syscall.Event
	for _, e := range events {
		sum |= syscall.AtomicLoadEvent(&e.raw)
	}
	sum &^= 1
	var ret uint32
	for {
		for n, e := range events {
			state := syscall.AtomicLoadEvent(&e.raw)
			state0, state1 := state&^1, state|1
			if syscall.AtomicCompareAndSwapEvent(&e.raw, state0, state1) {
				ret |= 1 << uint(n)
			}
		}
		if ret != 0 || t >= 0 && Nanosec() >= t {
			return ret
		}
		syscall.SetAlarm(t)
		sum.Wait()
	}
}
