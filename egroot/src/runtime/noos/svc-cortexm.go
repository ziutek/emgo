// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"syscall"
	"unsafe"

	"arch/cortexm"
	"arch/cortexm/debug/itm"
	"arch/cortexm/nvic"
)

//emgo:const
var syscalls = [...]func(fp *cortexm.StackFrame, lr uintptr){
	syscall.NEWTASK:       scNewTask,
	syscall.KILLTASK:      scKillTask,
	syscall.TASKUNLOCK:    scTaskUnlock,
	syscall.MAXTASKS:      scMaxTasks,
	syscall.SCHEDNEXT:     scSchedNext,
	syscall.EVENTWAIT:     scEventWait,
	syscall.SETSYSTIM:     scSetSysTimer,
	syscall.NANOSEC:       scNanosec,
	syscall.SETALARM:      scSetAlarm,
	syscall.SETIRQENA:     scSetIRQEna,
	syscall.SETIRQPRIO:    scSetIRQPrio,
	syscall.SETIRQHANDLER: scSetIRQHandler,
	syscall.IRQSTATUS:     scIRQStatus,
	syscall.SETPRIVLEVEL:  scSetPrivLevel,
	syscall.DEBUGOUT:      scDebugOut,
}

func unpriv() bool {
	return cortexm.CONTROL()&cortexm.Unpriv != 0
}

func mustThread(lr uintptr) {
	if lr&cortexm.ReturnMask == cortexm.ReturnHandler {
		panic("syscall from ISR")
	}
}

func scNewTask(fp *cortexm.StackFrame, lr uintptr) {
	mustThread(lr)
	tid, err := tasker.newTask(fp.R[0], fp.PSR, fp.R[1] != 0)
	fp.R[0], fp.R[1] = uintptr(tid), uintptr(err)
}

func scKillTask(fp *cortexm.StackFrame, lr uintptr) {
	mustThread(lr)
	err := tasker.killTask(int(fp.R[0]))
	fp.R[1] = uintptr(err)
}

func scTaskUnlock(fp *cortexm.StackFrame, lr uintptr) {
	mustThread(lr)
	tasker.unlockParent()
}

func scMaxTasks(fp *cortexm.StackFrame, lr uintptr) {
	fp.R[0] = uintptr(maxTasks())
}

func scSchedNext(_ *cortexm.StackFrame, lr uintptr) {
	raisePendSV()
}

func scEventWait(fp *cortexm.StackFrame, lr uintptr) {
	mustThread(lr)
	tasker.waitEvent(syscall.Event(fp.R[0]))
}

func scSetSysTimer(fp *cortexm.StackFrame, lr uintptr) {
	mustThread(lr)
	tasker.SetNanosec(utofr64(fp.R[0]))
	tasker.SetWakeup(utof64(fp.R[1]))
	// Raise PendSV to cause to call wakeup for the first time. System timer
	// implementation can rely on the fact, that wakeup is called by PendSV
	// handler immediately after setup (can be used to initialize/start timer).
	raisePendSV() 
}

func scNanosec(fp *cortexm.StackFrame, lr uintptr) {
	*(*int64)(unsafe.Pointer(&fp.R[0])) = tasker.Nanosec()
}

func scSetAlarm(fp *cortexm.StackFrame, lr uintptr) {
	mustThread(lr)
	tasker.SetAlarm(*(*int64)(unsafe.Pointer(&fp.R[0])))
}

func scSetIRQEna(fp *cortexm.StackFrame, lr uintptr) {
	irq := nvic.IRQ(fp.R[0])
	ena := fp.R[1] != 0
	if irq > 239 {
		fp.R[1] = uintptr(syscall.ERANGE)
		return
	}
	if ena {
		irq.Enable()
	} else {
		irq.Disable()
	}
	fp.R[1] = uintptr(syscall.OK)
}

func scSetIRQPrio(fp *cortexm.StackFrame, lr uintptr) {
	irq := nvic.IRQ(fp.R[0])
	prio := int(fp.R[1])
	if irq > 239 {
		fp.R[1] = uintptr(syscall.ERANGE)
		return
	}
	irq.SetPrio(prio)
	fp.R[1] = uintptr(syscall.OK)
}

func scSetIRQHandler(fp *cortexm.StackFrame, lr uintptr) {
	fp.R[1] = uintptr(syscall.ENOSYS)
}

func scIRQStatus(fp *cortexm.StackFrame, lr uintptr) {
	irq := nvic.IRQ(fp.R[0])
	if irq > 239 {
		fp.R[1] = uintptr(syscall.ERANGE)
		return
	}
	status := uintptr(irq.Prio())
	if irq.Enabled() {
		status = -status - 1
	}
	fp.R[0] = uintptr(status)
	fp.R[1] = uintptr(syscall.OK)
}

func scSetPrivLevel(fp *cortexm.StackFrame, lr uintptr) {
	mustThread(lr)
	ctrl := cortexm.CONTROL()
	switch level := int(fp.R[0]); {
	case level == 0:
		cortexm.SetCONTROL(ctrl &^ cortexm.Unpriv)
	case level > 0:
		cortexm.SetCONTROL(ctrl | cortexm.Unpriv)
	}
	fp.R[0] = uintptr(ctrl & cortexm.Unpriv)
	fp.R[1] = uintptr(syscall.OK)
}

func scDebugOut(fp *cortexm.StackFrame, lr uintptr) {
	port := itm.Port(fp.R[0])
	data := (*[1 << 30]byte)(unsafe.Pointer(fp.R[1]))[:fp.R[2]:fp.R[2]]
	if unpriv() && port >= 16 {
		fp.R[0] = 0
		fp.R[1] = uintptr(syscall.EPERM)
		return
	}
	n, _ := port.Write(data)
	fp.R[0] = uintptr(n)
	fp.R[1] = uintptr(syscall.OK)
}

// svcHandler calls sv with SVC caller's stack frame.
func svcHandler()

// Consider pass syscall number as a parameter instead embed it into SVC
// instruction. It take me few hours to analyze a bug caused by software
// breakpoints: (fp.PC - 2) points to the number embeded in BKPT instruction
// (that was inserted by gdb) instead of number in SVC instruction, but gdb x
// command shows right value and the fun begins...
//
// Tried syscall number in r0.
// Pros: avoid above issue, syscal number can be variable, only one read from
// SRAM to obtain number (embeded number need additional read from Flash).
// Cons: additional register need for syscall number, usually additional
// mov instruction is need (+2B for any syscall), additional instruction fetch
// from Flash + execution.
func sv(fp *cortexm.StackFrame, lr uintptr) {
	trap := int(*(*byte)(unsafe.Pointer(fp.PC - 2)))
	if trap >= len(syscalls) {
		panic("unknown syscall number")
	}
	syscalls[trap](fp, lr)
}
