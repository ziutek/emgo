// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"syscall"

	"arch/cortexm"
	"arch/cortexm/scb"
)

func raisePendSV() {
	scb.ICSR_Store(scb.PENDSVSET)
}

// pendSVHandler calls nextTask with current task PSP. It performs
// context swich if nextTask returns new non-zero value for PSP.
func pendSVHandler()

// nextTask returns taskInfo.sp for next task or 0.
// TODO: better scheduler
func nextTask(sp uintptr) uintptr {
again:
	t := uptime()
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
			tasker.setWakeup(tasker.alarm)
			cortexm.WFE()
			goto again
		}
	}
	if t += int64(tasker.period); t > tasker.alarm {
		t = tasker.alarm
	}
	tasker.setWakeup(t)
	if n == tasker.curTask {
		return 0
	}
	tasker.tasks[tasker.curTask].sp = sp
	tasker.curTask = n
	return tasker.tasks[n].sp
}
