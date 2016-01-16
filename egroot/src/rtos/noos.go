// +build noos

package rtos

import "syscall"

func sleepUntil(end int64) {
	for Nanosec() < end {
		syscall.SetAlarm(end)
		syscall.Alarm.Wait()
	}
}
