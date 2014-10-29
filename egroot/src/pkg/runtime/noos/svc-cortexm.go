// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"unsafe"

	"cortexm/exce"
)

// svcHandler calls sv with SVC caller's stack frame.
func svcHandler()

func sv(fp *exce.StackFrame) {
	// Consider pass SVC number as a parameter instead embed it into SVC
	// instruction. It take me few hours to analyze a bug caused by software
	// breakpoints: the following line returns number embeded in BKPT
	// instruction (that was inserted by gdb) instead of number in SVC
	// instruction, but gdb x command shows right values and the fun begins...
	n := *(*byte)(unsafe.Pointer(fp.PC - 2))
	switch n {
	case 0:
		tasker.newTask(fp.R[0], fp.PSR, fp.R[1] != 0)

	case 1:
		tasker.delTask(tasker.curTask)

	case 2:
		tasker.run()

	case 3:
		tasker.waitEvent(Event(fp.R[0]))
	}
}

func (ts *taskSched) newTask(pc uintptr, psr uint32, wait bool) {
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

	sf, sp := allocStackFrame(stackTop(n))
	ts.tasks[n].init()
	ts.tasks[n].sp = sp

	// Use parent's PSR as initial PSR for new task.
	sf.PSR = psr
	sf.PC = pc

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
