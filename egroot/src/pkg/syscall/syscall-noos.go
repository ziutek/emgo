// build +noos

package syscall

import "builtin"

const (
	NEWTASK    = builtin.NEWTASK
	DELTASK    = builtin.DELTASK
	TASKUNLOCK = builtin.TASKUNLOCK
	EVENTWAIT  = iota
	SETSYSCLK
	UPTIME
	SETIRQENA
	SETIRQPRIO
	SETISR
	IRQSTATUS
)

// NewTask creates new task that starts execute f. If lock is true tasker stops
// scheduling current task and waits until new task will call TaskUnlock. When
// success it returns TID of new task.
func NewTask(f func(), lock bool) (int, Errno) {
	tid, err := builtin.Syscall2(NEWTASK, f2p(f), b2p(lock))
	return int(tid), Errno(err)
}

// DelTask kills task with specified tid. tid == 0 means current task.
func DelTask(tid int) Errno {
	_, err := builtin.Syscall1(DELTASK, uintptr(tid))
	return Errno(err)
}

// TaskUnlock can be used when task was created with lock option. It informs
// tasker that now it can safely run parent task.
func TaskUnlock() {
	builtin.Syscall0(TASKUNLOCK)
}

// EventWait waits for event e.
func EventWait(e uintptr) {
	builtin.Syscall1(EVENTWAIT, e)
}

// SetSysClock informs runtime about current system clock frequency.
// It should be called at every system clock change.
func SetSysClock(hz uint) Errno {
	_, err := builtin.Syscall1(SETSYSCLK, uintptr(hz))
	return Errno(err)
}

// Uptime returns how long system is running (in nanosecond). Time when system
// was in deep sleep state can be included or not.
func Uptime() uint64 {
	return builtin.Syscall0u64(UPTIME)
}

// SetIRQEna enables or disables irq.
func SetIRQEna(irq int, ena bool) Errno {
	_, err := builtin.Syscall1(SETIRQENA, b2p(ena))
	return Errno(err)
}

// SetIRQPrio sets priority for irq.
func SetIRQPrio(irq, prio int)