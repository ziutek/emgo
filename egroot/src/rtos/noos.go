// +build noos

package rtos

import (
	"sync/fence"
	"syscall"
)

func sleepUntil(end int64) {
	for Nanosec() < end {
		syscall.SetAlarm(end)
		syscall.Alarm.Wait()
	}
}

func at(t int64) <-chan int64 {
	syscall.SetAt(t)
	return syscall.TimeChan()
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

func (f *EventFlag) reset(val int) {
	fence.RW_SMP() // Reset has RELEASE semantic.
	state := atomicInitLoadState(&f.state)
	event := state &^ 1
	new := event | syscall.Event(val&1)
	if state != new {
		syscall.AtomicStoreEvent(&f.state, new)
	}
}

func (f *EventFlag) signal(val int) {
	fence.RW_SMP() // Signal has RELEASE semantic.
	state := atomicInitLoadState(&f.state)
	event := state &^ 1
	new := event | syscall.Event(val&1)
	if state != new {
		syscall.AtomicStoreEvent(&f.state, new)
		event.Send()
	}
}

func (f *EventFlag) value() int {
	// Not need  atomicInitLoadState because an uninitialized f is zero.
	v := int(syscall.AtomicLoadEvent(&f.state) & 1)
	fence.RW_SMP()
	return v
}

func (f *EventFlag) wait(val int, deadline int64) (done bool) {
	state := atomicInitLoadState(&f.state)
	event := state &^ 1
	need := event | syscall.Event(val&1)
	if deadline != 0 {
		event |= syscall.Alarm
	}
	for {
		done = (state == need)
		if done {
			break
		}
		if deadline != 0 {
			if Nanosec() >= deadline {
				break
			}
			syscall.SetAlarm(deadline)
		}
		event.Wait()
		state = syscall.AtomicLoadEvent(&f.state)
	}
	fence.RW_SMP() // Wait has ACQUIRE semantic.
	return
}

func waitEvent(val int, deadline int64, flags []*EventFlag) (ret uint32) {
	var sum syscall.Event
	if deadline != 0 {
		sum = syscall.Alarm
	}
	for n, f := range flags {
		state := atomicInitLoadState(&f.state)
		sum |= state
		ret |= uint32(state&1) << uint(n)
	}
	sum &^= 1
	var v32 uint32
	if val&1 == 0 {
		v32 = 1<<uint(len(flags)) - 1
	}
	for {
		if ret != v32 {
			break
		}
		if deadline != 0 {
			if Nanosec() >= deadline {
				break
			}
			syscall.SetAlarm(deadline)
		}
		sum.Wait()
		ret = 0
		for n, f := range flags {
			state := syscall.AtomicLoadEvent(&f.state)
			ret |= uint32(state&1) << uint(n)
		}
	}
	fence.RW_SMP() // WaitEvent has ACQUIRE semantic.
	return
}
