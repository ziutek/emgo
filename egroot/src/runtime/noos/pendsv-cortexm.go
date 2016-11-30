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
// TODO: better scheduler

func nextTask(sp uintptr) uintptr {
	if softStackGuard {
		checkStackGuard(tasker.curTask)
	}
again:
	if ereg := syscall.TakeEventReg(); ereg != 0 {
		if ereg&syscall.Alarm != 0 {
			tasker.alarm = noalarm
		}
		tasker.deliverEvent(ereg)
	}
	n := tasker.curTask
	for {
		if n++; n == len(tasker.tasks) {
			n = 0
		}
		if tasker.tasks[n].state() == taskReady {
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
	tasker.tasks[tasker.curTask].sp = sp
	tasker.curTask = n
	wkup := tasker.nanosec() + int64(tasker.period)
	if wkup > tasker.alarm {
		wkup = tasker.alarm
	}
	if useMPU {
		setMPUStackGuard(n)
	}
	tasker.setWakeup(wkup, false)
	return tasker.tasks[n].sp
}
