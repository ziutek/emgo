// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"unsafe"

	"cortexm/exce"
)

// svcHandler calls sv with SVC caller's stack frame.
func svcHandler()

func sv(fp *stackFrame) {
	// Consider pass SVC number as a parameter instead embed it into SVC
	// instruction. It take me few hours to analyze a bug caused by software
	// breakpoints: the following line returns number embeded in BKPT
	// instruction (that was inserted by gdb) instead of number in SVC
	// instruction, but x command shows right values and the fun begins...
	n := *(*byte)(unsafe.Pointer(fp.pc - 2))
	switch n {
	case 0:
		tasker.newTask(fp.r[0], fp.xpsr, fp.r[1] != 0)

	case 1:
		tasker.delTask(tasker.curTask)

	case 2:
		tasker.run()

	case 3:
		tasker.waitEvent(Event(fp.r[0]))
	}
}

func (ts *taskSched) newTask(pc uintptr, xpsr uint32, wait bool) {
	n := ts.curTask
	for {
		if n++; n >= len(ts.tasks) {
			n = 0
		}
		if ts.tasks[n].state() == taskEmpty {
			break
		}
		if n == ts.curTask {
			panic("too many tasks")
		}
	}

	sf, sp := allocStackFrame(initSP(n))
	ti := taskInfo{sp: sp, prio: 255}
	ti.rng.Seed(Ticks())
	ts.tasks[n] = ti

	// Use parent's xPSR as initial xPSR for new task.
	sf.xpsr = xpsr
	sf.pc = pc

	ts.tasks[n].setState(taskReady)

	if wait {
		ts.stop()
		// BUG: This badly affects scheduling but don't care for now.
		ts.forceNext = n
		exce.PendSV.SetPending()
	}
}

func (ts *taskSched) delTask(n int) {
	ts.tasks[n].setState(taskEmpty)
	exce.PendSV.SetPending()
}

func (ts *taskSched) waitEvent(e Event) {
	t := &ts.tasks[ts.curTask]
	if e == 0 || t.event&e != 0 {
		t.event = 0
		return
	}
	t.setState(taskWaitEvent)
	t.event = e
	exce.PendSV.SetPending()
}
