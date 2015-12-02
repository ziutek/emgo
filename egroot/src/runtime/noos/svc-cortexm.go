// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm"
	"arch/cortexm/debug/itm"
	"arch/cortexm/nvic"
	"syscall"
	"unsafe"
)

var syscalls = [...]func(*cortexm.StackFrame){
	syscall.NEWTASK:       scNewTask,
	syscall.KILLTASK:      scKillTask,
	syscall.TASKUNLOCK:    scTaskUnlock,
	syscall.SCHEDNEXT:     scSchedNext,
	syscall.EVENTWAIT:     scEventWait,
	syscall.SETSYSCLK:     scSetSysClock,
	syscall.UPTIME:        scUptime,
	syscall.SETIRQENA:     scSetIRQEna,
	syscall.SETIRQPRIO:    scSetIRQPrio,
	syscall.SETIRQHANDLER: scSetIRQHandler,
	syscall.IRQSTATUS:     scIRQStatus,
	syscall.DEBUGOUT:      scDebugOut,
}

func unpriv() bool {
	return cortexm.CONTROL()&cortexm.Unpriv != 0
}

func scNewTask(fp *cortexm.StackFrame) {
	tid, err := tasker.newTask(fp.R[0], fp.PSR, fp.R[1] != 0)
	fp.R[0], fp.R[1] = uintptr(tid), uintptr(err)
}

func scKillTask(fp *cortexm.StackFrame) {
	err := tasker.killTask(int(fp.R[0]))
	fp.R[1] = uintptr(err)
}

func scTaskUnlock(_ *cortexm.StackFrame) {
	tasker.unlockParent()
}

func scSchedNext(_ *cortexm.StackFrame) {
	raisePendSV()
}

func scEventWait(fp *cortexm.StackFrame) {
	tasker.waitEvent(syscall.Event(fp.R[0]))
}

func scSetSysClock(fp *cortexm.StackFrame) {
	sysTimerHz = uint64(fp.R[0])
	setTimerFreq(uint32(fp.R[0]))
}

func scUptime(fp *cortexm.StackFrame) {
	*(*uint64)(unsafe.Pointer(fp)) = uptime()
}

func scSetIRQEna(fp *cortexm.StackFrame) {
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

func scSetIRQPrio(fp *cortexm.StackFrame) {
	irq := nvic.IRQ(fp.R[0])
	prio := int(fp.R[1])
	if irq > 239 {
		fp.R[1] = uintptr(syscall.ERANGE)
		return
	}
	irq.SetPrio(prio)
	fp.R[1] = uintptr(syscall.OK)
}

func scSetIRQHandler(fp *cortexm.StackFrame) {
	fp.R[1] = uintptr(syscall.ENOSYS)
}

func scIRQStatus(fp *cortexm.StackFrame) {
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

func scDebugOut(fp *cortexm.StackFrame) {
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
func sv(fp *cortexm.StackFrame) {
	trap := int(*(*byte)(unsafe.Pointer(fp.PC - 2)))
	if trap >= len(syscalls) {
		panic("unknown syscall number")
	}
	syscalls[trap](fp)
}
