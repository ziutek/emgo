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

type eventFlag struct {
	state syscall.Event
}

func atomicInitLoadState(p *syscall.Event) syscall.Event {
	state := syscall.AtomicLoadEvent(p)
	if state == 0 {
		state = syscall.AssignEventFlag() | syscall.Alarm
		if !syscall.AtomicCompareAndSwapEvent(p, 0, state) {
			state = syscall.AtomicLoadEvent(p)
		}
	}
	return state
}

func flagSet(f *EventFlag) {
	state := atomicInitLoadState(&f.state) | 1
	syscall.AtomicStoreEvent(&f.state, state)
	state.Send()
}

func flagClear(f *EventFlag) {
	state := atomicInitLoadState(&f.state) &^ 1
	syscall.AtomicStoreEvent(&f.state, state)
	state.Send()
}

func flagVal(f *EventFlag) int {
	return int(syscall.AtomicLoadEvent(&f.state) & 1)
}

func flagWait(f *EventFlag, deadline int64) bool {
	for {
		state := atomicInitLoadState(&f.state)
		if state&1 != 0 {
			return true
		}
		if deadline != 0 && Nanosec() >= deadline {
			return false
		}
		syscall.SetAlarm(deadline)
		(state &^ 1).Wait()
	}
}

func waitEvent(deadline int64, flags []*EventFlag) uint32 {
	var sum syscall.Event
	for _, f := range flags {
		sum |= atomicInitLoadState(&f.state)
	}
	sum &^= 1
	var ret uint32
	for {
		for n, f := range flags {
			state := syscall.AtomicLoadEvent(&f.state)
			ret |= uint32(state&1) << uint(n)
		}
		if ret != 0 || deadline != 0 && Nanosec() >= deadline {
			return ret
		}
		syscall.SetAlarm(deadline)
		sum.Wait()
	}
}
