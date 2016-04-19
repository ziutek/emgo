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
		state = syscall.AssignEventFlag()
		if !syscall.AtomicCompareAndSwapEvent(p, 0, state) {
			state = syscall.AtomicLoadEvent(p)
		}
	}
	return state
}

func flagSet(f *EventFlag) {
	state := atomicInitLoadState(&f.state) | 1
	syscall.AtomicStoreEvent(&f.state, state)
	(state &^ 1).Send()
}

func flagClear(f *EventFlag) {
	state := atomicInitLoadState(&f.state) &^ 1
	syscall.AtomicStoreEvent(&f.state, state)
}

func flagVal(f *EventFlag) int {
	// This code rely on the fact that uninitialized f is zero.
	return int(syscall.AtomicLoadEvent(&f.state) & 1)
}

func flagWait(f *EventFlag, deadline int64) bool {
	state := atomicInitLoadState(&f.state)
	if deadline != 0 {
		state |= syscall.Alarm
	}
	for {
		if state&1 != 0 {
			return true
		}
		if deadline != 0 {
			if Nanosec() >= deadline {
				return false
			}
			syscall.SetAlarm(deadline)
		}
		state.Wait()
		state = syscall.AtomicLoadEvent(&f.state)
	}
}

func waitEvent(deadline int64, flags []*EventFlag) uint32 {
	var (
		sum syscall.Event
		ret uint32
	)
	if deadline != 0 {
		sum = syscall.Alarm
	}
	for n, f := range flags {
		state := atomicInitLoadState(&f.state)
		sum |= state
		ret |= uint32(state&1) << uint(n)
	}
	sum &^= 1
	for {
		if ret != 0 {
			return ret
		}
		if deadline != 0 {
			if Nanosec() >= deadline {
				return 0
			}
			syscall.SetAlarm(deadline)
		}
		sum.Wait()
		for n, f := range flags {
			state := syscall.AtomicLoadEvent(&f.state)
			ret |= uint32(state&1) << uint(n)
		}
	}
}
