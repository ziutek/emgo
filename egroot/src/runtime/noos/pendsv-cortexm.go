// +build cortexm0 cortexm3 cortexm4 cortexm4f

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
	ns := tasker.nanosec()
	if ns >= tasker.alarm {
		syscall.Alarm.Send()
		tasker.alarm = 1<<63 - 1
	}
	if ereg := syscall.TakeEventReg(); ereg != 0 {
		tasker.deliverEvent(syscall.Event(ereg))
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
			//dbg.WriteString("*nt* 0\n")
			tasker.wakeup(tasker.alarm)
			cortexm.WFE()
			//dbg.WriteString("*nt* wakeup\n")
			clearPendSV() // Avoid reentering after return.
			goto again
		}
	}
	if n == tasker.curTask {
		// Only one task running.
		//dbg.WriteString("*nt* 1\n")
		tasker.wakeup(tasker.alarm)
		return 0
	}
	tasker.tasks[tasker.curTask].sp = sp
	tasker.curTask = n
	if ns += int64(tasker.period); ns > tasker.alarm {
		ns = tasker.alarm
	}
	//dbg.WriteString("*nt* N\n")
	if hasMPU {
		setMPUStackGuard(n)
	}
	tasker.wakeup(ns)
	return tasker.tasks[n].sp
}
