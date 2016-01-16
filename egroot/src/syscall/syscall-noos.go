// build +noos

package syscall

import (
	"bits"
	"builtin"
	"unsafe"
)

const (
	NEWTASK    = builtin.NEWTASK
	KILLTASK   = builtin.KILLTASK
	TASKUNLOCK = builtin.TASKUNLOCK
	MAXTASKS   = iota
	SCHEDNEXT
	EVENTWAIT
	SETSYSCLK
	NANOSEC
	SETALARM
	SETIRQENA
	SETIRQPRIO
	SETIRQHANDLER
	IRQSTATUS
	SETPRIVLEVEL
	DEBUGOUT
)

// NewTask creates new task that starts execute f. If lock is true tasker stops
// scheduling current task and waits until new task will call TaskUnlock. When
// success it returns TID of new task.
func NewTask(f func(), lock bool) (int, Errno) {
	tid, e := builtin.Syscall2(NEWTASK, ftou(f), uintptr(bits.One(lock)))
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

// MaxTasks: see rtos package.
func MaxTasks() int {
	n, _ := builtin.Syscall0(MAXTASKS)
	return int(n)
}

// SchedNext informs tasker that it need to schedule next ready to run task.
// It is safe to call SchedNext from interrupt handler.
func SchedNext() {
	schedNext()
}

// SetSysClock registers two functions that runtime uses to communicate with
// system clock.
//
// nanosec is used to implement Nanosec system call. It should return the
// monotonic time i nanoseconds (typically the time of system clock run).
//
// wakeup is called by scheduler to ask system clock to generate PendSV
// exception at time t (see rtos.Nanosec). Weak (ticking) system clock can
// ignore t and generate PendSV with fixed period. Good (tickless) clock should
// generate exceptions as accurately as possible (if t <= nanosec() it should
// generate PendSV immediately. wakeup must not generate any exception before
// return, to do not wakeup runtime too early from WFE.
func SetSysClock(nanosec func() int64, wakeup func(t int64)) Errno {
	_, e := builtin.Syscall2(SETSYSCLK, fr64tou(nanosec), f64tou(wakeup))
	return Errno(e)
}

// Nanosec: see rtos package.
func Nanosec() int64 {
	return builtin.Syscall0r64(NANOSEC)
}

// SetAlarm asks the runtime to send Alarm event at t. t is
// only a hint for runtime, because it can send alarm at any time: before t,
// at t and after t. Typically, task use SetAlarm in conjunction with
// Alarm.Wait and Nanosec.
func SetAlarm(t int64) {
	builtin.Syscall1i64(SETALARM, t)
}

// SetIRQEna enables or disables irq.
func SetIRQEna(irq int, ena bool) Errno {
	_, e := builtin.Syscall2(SETIRQENA, uintptr(irq), uintptr(bits.One(ena)))
	return Errno(e)
}

// SetIRQPrio sets priority for irq.
func SetIRQPrio(irq, prio int) Errno {
	_, err := builtin.Syscall2(SETIRQPRIO, uintptr(irq), uintptr(prio))
	return Errno(err)
}

// SetIRQHandler: see rtos package.
func SetIRQHandler(irq int, f func()) Errno {
	_, e := builtin.Syscall2(SETIRQHANDLER, uintptr(irq), ftou(f))
	return Errno(e)
}

// IRQStatus: ee rtos package.
func IRQStatus(irq int) (int, Errno) {
	s, e := builtin.Syscall1(IRQSTATUS, uintptr(irq))
	return int(s), Errno(e)
}

// SetPrivLevel: see rtos package.
func SetPrivLevel(n int) (int, Errno) {
	old, e := builtin.Syscall1(SETPRIVLEVEL, uintptr(n))
	return int(old), Errno(e)
}

// DebugOutString allows write debug message.
func DebugOutString(port int, s string) (int, Errno) {
	p := (*builtin.String)(unsafe.Pointer(&s))
	n, e := builtin.Syscall3(DEBUGOUT, uintptr(port), p.Addr, p.Len)
	return int(n), Errno(e)
}

// DebugOut allows write debug message.
func DebugOut(port int, b []byte) (int, Errno) {
	return DebugOutString(port, *(*string)(unsafe.Pointer(&b)))
}
