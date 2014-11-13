// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"cortexm/exce"
	"syscall"
	"unsafe"
)

var syscalls = [...]func(*exce.StackFrame){
	syscall.NEWTASK:    newTask,
	syscall.DELTASK:    delTask,
	syscall.TASKUNLOCK: taskUnlock,
	syscall.EVENTWAIT:  eventWait,
	syscall.SETSYSCLK:  setSysClock,
}

func newTask(fp *exce.StackFrame) {
	tid, err := tasker.newTask(fp.R[0], fp.PSR, fp.R[1] != 0)
	fp.R[0], fp.R[1] = uintptr(tid), uintptr(err)
}

func delTask(fp *exce.StackFrame) {
	err := tasker.delTask(int(fp.R[0]))
	fp.R[0] = uintptr(err)
}

func taskUnlock(fp *exce.StackFrame) {
	tasker.unlockParent()
}

func eventWait(fp *exce.StackFrame) {
	tasker.waitEvent(Event(fp.R[0]))
}

func setSysClock(fp *exce.StackFrame) {
	sysClk = uint(fp.R[0])
	setTickPeriod()
}

// svcHandler calls sv with SVC caller's stack frame.
func svcHandler()

// Consider pass syscall number as a parameter instead embed it into SVC
// instruction. It take me few hours to analyze a bug caused by software
// breakpoints: the following line returns number embeded in BKPT
// instruction (that was inserted by gdb) instead of number in SVC
// instruction, but gdb x command shows right values and the fun begins...
//
// Tried syscall number in r0.
// Pros: avoid above issue, syscal number can be variable, only one read from
// SRAM to obtain number (embeded number need additional read from Flash).
// Cons: additional register need for syscall number, usually additional
// mov instruction is need (+2B for any syscall), additional instruction fetch
// from Flash + execution.

func sv(fp *exce.StackFrame) {
	trap := int(*(*byte)(unsafe.Pointer(fp.PC - 2)))
	if trap >= len(syscalls) {
		panic("unknown syscall number")
	}
	syscalls[trap](fp)
}
