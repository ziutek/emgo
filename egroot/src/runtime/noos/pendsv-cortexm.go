// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package noos

import (
	"syscall"

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
	now := tasker.nanosec()
	clearPendSV() // Some ISR could set it again. Clear just before takeEvents.
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
			now = tasker.nanosec()
			clearPendSV() // Clear PENDSV flag just before takeEvents.
			if ev := tasker.takeEvents(now); ev != 0 {
				tasker.deliverEvents(ev)
				break
			}
			cortexm.WFE()
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
	if useMPU {
		setMPUStackGuard(nextTask)
	}
	tasker.setWakeup(wkup)
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
*/
