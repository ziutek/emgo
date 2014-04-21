// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import "sync/atomic"

// pendSVHandler calls nextTask with PSP for current task. It performs
// context swich if nextTask returns new non-zero value for PSP.
func pendSVHandler()

// nextTask returns taskInfo.sp for next task or 0.
// TODO: better scheduler
func nextTask(sp uintptr) uintptr {
	if ereg := atomic.SwapUintptr((*uintptr)(&eventReg), 0); ereg != 0 {
		tasker.deliverEvent(Event(ereg))
	}

	n := tasker.forceNext
	if n >= 0 {
		tasker.forceNext = -1
	} else {
		n = tasker.curTask
		for {
			if n++; n >= len(tasker.tasks) {
				n = 0
			}
			if tasker.tasks[n].state() == taskReady {
				break
			}
			if n == tasker.curTask {
				panic("no task to run")
			}
		}
		if n == tasker.curTask {
			return 0
		}
	}
	tasker.tasks[tasker.curTask].sp = sp
	tasker.curTask = n
	return tasker.tasks[n].sp
}
