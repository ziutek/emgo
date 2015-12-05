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
	SCHEDNEXT  = iota
	EVENTWAIT
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
	tid, e := builtin.Syscall2(NEWTASK, ftou(f), btou(lock))
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

// SchedNext informs tasker that it need to schedule next ready to run task.
// It is safe to call SchedNext from interrupt handler.
func SchedNext() {
	schedNext()
}

// SetSysClock informs runtime about current system clock frequency. If rtc
// is nil generic clock source is used (eg. Cortex-M SysTick timer). Otherwise
// rtc is called to obtain Real Time Counter ticks since last update of
// sysclock. Use SetSysClock(0, nil) to disable system clock (default state
// after reset).
func SetSysClock(hz uint, rtc func() uint32) Errno {
	_, e := builtin.Syscall2(SETSYSCLK, uintptr(hz), f32tou(rtc))
	return Errno(e)
}

// Uptime returns how long system is running (in nanosecond). Time when system
// was in deep sleep state can be included or not.
func Uptime() uint64 {
	return builtin.Syscall0u64(UPTIME)
}

// SetIRQEna enables or disables irq.
func SetIRQEna(irq int, ena bool) Errno {
	_, e := builtin.Syscall2(SETIRQENA, uintptr(irq), btou(ena))
	return Errno(e)
}

// SetIRQPrio sets priority for irq.
func SetIRQPrio(irq, prio int) Errno {
	_, err := builtin.Syscall2(SETIRQPRIO, uintptr(irq), uintptr(prio))
	return Errno(err)
}

// SetIRQHandler sets f as handler function for irq.
func SetIRQHandler(irq int, f func()) Errno {
	_, e := builtin.Syscall2(SETIRQHANDLER, uintptr(irq), ftou(f))
	return Errno(e)
}

func IRQStatus(irq int) (int, Errno) {
	s, e := builtin.Syscall1(IRQSTATUS, uintptr(irq))
	return int(s), Errno(e)
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

// SetAlarm asks the runtime to send Alarm event a uptime t nanoseconds. t is
// only a hint for runtime, because it can send alarm at any time: before t,
// at t and after t. Typically, task use SetAlarm in conjunction with
// Alarm.Wait and Uptime.
func SetAlarm(t uint64) {

}