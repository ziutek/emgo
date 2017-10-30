// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package noos

import (
	"syscall"
	"sync/fence"

	"arch/cortexm"
	"arch/cortexm/scb"
)

//const dbg = itm.Port(17)

func raisePendSV() { scb.SCB.ICSR.Store(scb.PENDSVSET) }
func clearPendSV() { scb.SCB.ICSR.Store(scb.PENDSVCLR) }

// pendSVHandler calls nextTask with current task PSP. It performs
// context swich if nextTask returns new non-zero value for PSP.
func pendSVHandler()

// nextTask returns taskInfo.sp for next task or 0.
func nextTask(sp uintptr) uintptr {
	if softStackGuard {
		checkStackGuard(tasker.curTask)
	}
	clearPendSV() // Some ISR could set it again. Can clear before nanosec.
	fence.RW()
	now := tasker.nanosec()
	if ev := tasker.takeEvents(now); ev != 0 {
		tasker.deliverEvents(ev)
	}
	var nextTask int
	for {
		nextTask = tasker.selectTask()
		if nextTask >= 0 {
			break
		}
		// No task to run.
		tasker.setWakeup(tasker.alarm)
		for {
			cortexm.WFE()
			clearPendSV()
			fence.RW()
			now = tasker.nanosec()
			if ev := tasker.takeEvents(now); ev != 0 {
				tasker.deliverEvents(ev)
				break
			}
		}
	}
	if nextTask == tasker.curTask {
		// Only one task is running.
		tasker.setWakeup(tasker.alarm)
		return 0
	}
	tasker.tasks[tasker.curTask].info.sp = sp
	tasker.curTask = nextTask
	wkup := now + int64(tasker.period)
	if wkup > tasker.alarm {
		wkup = tasker.alarm
	}
	tasker.setWakeup(wkup)
	if useMPU {
		setMPUStackGuard(nextTask)
	}
	return tasker.tasks[nextTask].info.sp
}

func (ts *taskSched) takeEvents(now int64) syscall.Event {
	ev := syscall.TakeEventReg()
	if now >= tasker.alarm {
		ev |= syscall.Alarm
		tasker.alarm = maxAlarm
	}
	return ev
}

func (ts *taskSched) selectTask() int {
	n := tasker.curTask
	for {
		if n++; n == len(tasker.tasks) {
			n = 0
		}
		if tasker.tasks[n].info.state() == taskReady {
			return n
		}
		if n == tasker.curTask {
			return -1
		}
	}
}

/*
func nextTask(sp uintptr) uintptr {
	if softStackGuard {
		checkStackGuard(tasker.curTask)
	}
again:
	if ereg := syscall.TakeEventReg(); ereg != 0 {
		if ereg&syscall.Alarm != 0 {
			tasker.alarm = noalarm
		}
		tasker.deliverEvents(ereg)
	}
	n := tasker.curTask
	for {
		if n++; n == len(tasker.tasks) {
			n = 0
		}
		if tasker.tasks[n].info.state() == taskReady {
			break
		}
		if n == tasker.curTask {
			// No task to run.
			tasker.setWakeup(tasker.alarm, tasker.alarm != noalarm)
			cortexm.WFE()
			clearPendSV() // Avoid reentering after return.
			goto again
		}
	}
	if n == tasker.curTask {
		// Only one task running.
		tasker.setWakeup(tasker.alarm, tasker.alarm != noalarm)
		return 0
	}
	tasker.tasks[tasker.curTask].info.sp = sp
	tasker.curTask = n
	wkup := tasker.nanosec() + int64(tasker.period)
	if wkup > tasker.alarm {
		wkup = tasker.alarm
	}
	if useMPU {
		setMPUStackGuard(n)
	}
	tasker.setWakeup(wkup, false)
	return tasker.tasks[n].info.sp
}

// Debuging LEDs on Port103R.

//c:volatile
type port struct {
	cr   [2]uint32
	idr  uint16
	_    uint16
	odr  uint16
	_    uint16
	bsrr uint32
	brr  uint32
	lckr uint32
}

func (p *port) Set(n uint) {
	p.bsrr = 1 << n
	//delay.Loop(1e6)
}

func (p *port) Clear(n uint) {
	p.bsrr = 1 << (n + 16)
	//delay.Loop(1e6)
}

var pb = (*port)(unsafe.Pointer(uintptr(0x40010C00)))

const (
	led1 = 7
	led2 = 6
	led3 = 5
)

*/
