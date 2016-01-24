// Package rtc implements tickless system clock using real time clock/counter.
// RTC system clock uses up to 20 bytes from first backup registers to preserve
// its state. 8 bytes are used to implement rtos.Nanosec. If SetTime function is
// used to set current calendar time, a further 12 bytes are used to preserve
// RTC start time.
package rtc

import (
	"rtos"
	"time"
)

// Stetup setups RTC as system clock using as clock source LSE. freqHz should
// be set to LSE frequency.
func Setup(freqHz uint) {
	setup(freqHz)
}

// SetTime sets current calendar time. It does not affect rtos.Nanosec(). Only
// time.Now() is affected.
func SetTime(t time.Time) {
	up := time.Duration(rtos.Nanosec())
	t = t.Add(-up)
	setStartTime(t)
}

// ISR sgould be set as irq.RTCAlarm interrupt handler.
func ISR() {
	isr()
}

// Status return clock status. ok informs whether clock is configured and works
// properly, set informs whether clock is set.
func Status() (ok, set bool) {
	return status()
}
