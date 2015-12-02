// +build noos

package rtos

import "syscall"

func sleepUntil(end uint64) {
	for Uptime() < end {
		// syscall.SetAlarm(end)
		syscall.Alarm.Wait()
	}
}
