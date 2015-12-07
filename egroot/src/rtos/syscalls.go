package rtos

import "syscall"

// Uptime returns how long system is running (in nanosecond). Time when system
// was in deep sleep state can be included or not.
func Uptime() int64 {
	return syscall.Uptime()
}

// SleepUntil sleeps task until Uptime() < end.
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
