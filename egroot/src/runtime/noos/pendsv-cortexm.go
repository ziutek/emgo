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

// pendSVHandler calls nextTask with PSP for current task. It performs
// context swich if nextTask returns new non-zero value for PSP.
func pendSVHandler()

// nextTask returns taskInfo.sp for next task or 0.
// TODO: better scheduler
func nextTask(sp uintptr) uintptr {
again:
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
			cortexm.WFE()
			goto again
		}
	}
	if n == tasker.curTask {
		return 0
	}
	tasker.tasks[tasker.curTask].sp = sp
	tasker.curTask = n
	return tasker.tasks[n].sp
}
