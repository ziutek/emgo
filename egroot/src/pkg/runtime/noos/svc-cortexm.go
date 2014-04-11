// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import "unsafe"

// svcHandler calls sv with SVC caller's caller stack frame.
func svcHandler()

func sv(fp *stackFrame) {
	n := *(*byte)(unsafe.Pointer(fp.pc - 2))
	switch n {
	case 0: // New task
		newTask(fp.r[0], fp.xpsr)

	case 1:
		delTask()
	}
}
