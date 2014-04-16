// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import "unsafe"

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
	}
}
