// +build noos

package rtos

import "syscall"

func sleepUntil(end int64) {
	for Uptime() < end {
		syscall.SetAlarm(end)
		syscall.Alarm.Wait()
	}
}
