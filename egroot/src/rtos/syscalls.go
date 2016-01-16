package rtos

import "syscall"

// Nanosec returns system time as number of nanosecond from some time in the
// past. There is guarantee that system time is monotonic, however accuracy
// and linearity is implementation dependent. Usually, only systems that use
// real time clock/counter as time source can provide linear system time. In
// other cases the time flow can be affected by any enter into low power sleep
// state.
func Nanosec() int64 {
	return syscall.Nanosec()
}

// SleepUntil sleeps task until Nanosec() < end.
func SleepUntil(end int64) {
	sleepUntil(end)
}

// SetPrivLevel sets privilege level for current task to n. Level 0 is the most
// privileged and allows access to all system resources. Any resource available
// to level n is also available to level 0..n. If n < 0 the privilege level is
// not changed. SetPrivLevel returns previous level number and error.
func SetPrivLevel(n int) (int, error) {
	old, e := syscall.SetPrivLevel(n)
	return old, mkerror(e)
}

// MaxTasks returns number of tasks that can exists simultaneously. MaxTasks
// returns 0 if there is no tasker enabled in runtime.
func MaxTasks() int {
	return syscall.MaxTasks()
}
