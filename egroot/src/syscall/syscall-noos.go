// build +noos

package syscall

import (
	"builtin"
	"unsafe"
)

const (
	NEWTASK    = builtin.NEWTASK
	KILLTASK   = builtin.KILLTASK
	TASKUNLOCK = builtin.TASKUNLOCK
	EVENTWAIT  = iota
	SETSYSCLK
	UPTIME
	SETIRQENA
	SETIRQPRIO
	SETIRQHANDLER
	IRQSTATUS
	DEBUGOUT
)

// NewTask creates new task that starts execute f. If lock is true tasker stops
// scheduling current task and waits until new task will call TaskUnlock. When
// success it returns TID of new task.
func NewTask(f func(), lock bool) (int, Errno) {
	tid, e := builtin.Syscall2(NEWTASK, f2u(f), b2u(lock))
	return int(tid), Errno(e)
}

// KillTask kills task with specified tid. tid == 0 means current task.
func KillTask(tid int) Errno {
	_, e := builtin.Syscall1(KILLTASK, uintptr(tid))
	return Errno(e)
}

// TaskUnlock can be used when task was created with lock option. It informs
// tasker that now it can safely run parent task.
func TaskUnlock() {
	builtin.Syscall0(TASKUNLOCK)
}

// SetSysClock informs runtime about current system clock frequency.
// It should be called at every system clock change.
func SetSysClock(hz uint) Errno {
	_, e := builtin.Syscall1(SETSYSCLK, uintptr(hz))
	return Errno(e)
}

// Uptime returns how long system is running (in nanosecond). Time when system
// was in deep sleep state can be included or not.
func Uptime() uint64 {
	return builtin.Syscall0u64(UPTIME)
}

// SetIRQEna enables or disables irq.
func SetIRQEna(irq int, ena bool) Errno {
	_, e := builtin.Syscall2(SETIRQENA, uintptr(irq), b2u(ena))
	return Errno(e)
}

// SetIRQPrio sets priority for irq.
func SetIRQPrio(irq, prio int) Errno {
	_, err := builtin.Syscall2(SETIRQPRIO, uintptr(irq), uintptr(prio))
	return Errno(err)
}

// SetIRQHandler sets f as handler function for irq.
func SetIRQHandler(irq int, f func()) Errno {
	_, e := builtin.Syscall2(SETIRQHANDLER, uintptr(irq), f2u(f))
	return Errno(e)
}

func IRQStatus(irq int) (int, Errno) {
	s, e := builtin.Syscall1(IRQSTATUS, uintptr(irq))
	return int(s), Errno(e)
}

// DebugOut allows write debug informations.
func DebugOut(port int, data []byte) (int, Errno) {
	n, e := builtin.Syscall3(
		DEBUGOUT,
		uintptr(port), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)),
	)
	return int(n), Errno(e)
}

// DebugOutString allows write debug message.
func DebugOutString(port int, s string) (int, Errno) {
	data := (*[]byte)(unsafe.Pointer(&s))
	return DebugOut(port, *data)
}
