// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"syscall"

	"arch/cortexm"
	"arch/cortexm/debug/itm"
	"arch/cortexm/scb"
)

const dbg = itm.Port(17)

func raisePendSV() { scb.ICSR_Store(scb.PENDSVSET) }
func clearPendSV() { scb.ICSR_Store(scb.PENDSVCLR) }

// pendSVHandler calls nextTask with current task PSP. It performs
// context swich if nextTask returns new non-zero value for PSP.
func pendSVHandler()

// nextTask returns taskInfo.sp for next task or 0.
// TODO: better scheduler
func nextTask(sp uintptr) uintptr {
again:
	t := tasker.uptime()
	if t >= tasker.alarm {
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
		if tasker.tasks[n].State() == taskReady {
			break
		}
		if n == tasker.curTask {
			// No task to run.
			dbg.WriteString("nextTask 0\n")
			tasker.wakeup(tasker.alarm)
			cortexm.WFE()
			clearPendSV() // Avoid reentering after return.
			goto again
		}
	}
	if n == tasker.curTask {
		// Only one task running.
		dbg.WriteString("nextTask 1\n")
		tasker.wakeup(tasker.alarm)
		return 0
	}
	tasker.tasks[tasker.curTask].sp = sp
	tasker.curTask = n
	if t += int64(tasker.period); t > tasker.alarm {
		t = tasker.alarm
	}
	dbg.WriteString("nextTask N\n")
	tasker.wakeup(t)
	return tasker.tasks[n].sp
}
