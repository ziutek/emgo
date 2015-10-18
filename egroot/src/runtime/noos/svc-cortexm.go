// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm"
	"arch/cortexm/debug/itm"
	"arch/cortexm/exce"
	"syscall"
	"unsafe"
)

var syscalls = [...]func(*exce.StackFrame){
	syscall.NEWTASK:       scNewTask,
	syscall.KILLTASK:      scKillTask,
	syscall.TASKUNLOCK:    scTaskUnlock,
	syscall.EVENTWAIT:     scEventWait,
	syscall.SETSYSCLK:     scSetSysClock,
	syscall.UPTIME:        scUptime,
	syscall.SETIRQENA:     scSetIRQEna,
	syscall.SETIRQPRIO:    scSetIRQPrio,
	syscall.SETIRQHANDLER: scSetIRQHandler,
	syscall.IRQSTATUS:     nil,
	syscall.DEBUGOUT:      scDebugOut,
}

func unpriv() bool {
	return cortexm.Ctrl()&cortexm.Unpriv != 0
}

func scNewTask(fp *exce.StackFrame) {
	tid, err := tasker.newTask(fp.R[0], fp.PSR, fp.R[1] != 0)
	fp.R[0], fp.R[1] = uintptr(tid), uintptr(err)
}

func scKillTask(fp *exce.StackFrame) {
	err := tasker.killTask(int(fp.R[0]))
	fp.R[1] = uintptr(err)
}

func scTaskUnlock(fp *exce.StackFrame) {
	tasker.unlockParent()
}

func scEventWait(fp *exce.StackFrame) {
	tasker.waitEvent(syscall.Event(fp.R[0]))
}

func scSetSysClock(fp *exce.StackFrame) {
	sysClk = uint(fp.R[0])
	setTickPeriod()
}

func scUptime(fp *exce.StackFrame) {
	*(*uint64)(unsafe.Pointer(fp)) = uptime()
}

func scSetIRQEna(fp *exce.StackFrame) {
	irq := exce.Exce(fp.R[0])
	ena := fp.R[1] != 0
	if irq < exce.IRQ0 && unpriv() {
		fp.R[1] = uintptr(syscall.EPERM)
		return
	}
	if ena {
		irq.Enable()
	} else {
		irq.Disable()
	}
	fp.R[1] = uintptr(syscall.OK)
}

func scSetIRQPrio(fp *exce.StackFrame) {
	irq := exce.Exce(fp.R[0])
	prio := exce.Prio(fp.R[1])
	if irq < exce.IRQ0 && unpriv() {
		fp.R[1] = uintptr(syscall.EPERM)
		return
	}
	irq.SetPrio(prio)
	fp.R[1] = uintptr(syscall.OK)
}

func scSetIRQHandler(fp *exce.StackFrame) {
	irq := exce.Exce(fp.R[0])
	h := p2f(fp.R[1])
	if irq < exce.IRQ0 && unpriv() {
		fp.R[1] = uintptr(syscall.EPERM)
		return
	}
	irq.UseHandler(h)
	fp.R[1] = uintptr(syscall.OK)
}

func scDebugOut(fp *exce.StackFrame) {
	port := fp.R[0]
	data := (*[1 << 30]byte)(unsafe.Pointer(fp.R[1]))[:fp.R[2]:fp.R[2]]
	if unpriv() && port > 15 {
		fp.R[0] = 0
		fp.R[1] = uintptr(syscall.EPERM)
		return
	}
	n, _ := itm.Port(port).Write(data)
	fp.R[0] = uintptr(n)
	fp.R[1] = uintptr(syscall.OK)
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
