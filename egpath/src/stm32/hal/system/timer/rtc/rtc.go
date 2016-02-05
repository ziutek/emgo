// Package rtc implements tickless system timer using real time clock/counter.
// RTC system timer uses up to 20 bytes from first backup registers to preserve
// its state. 8 bytes are used to implement rtos.Nanosec. If SetTime function
// is used to set current calendar time, a further 12 bytes  are used to
// preserve RTC start time.
package rtc

import (
	"rtos"
	"time"
)

// Stetup setups RTC as system timer using LSE as clock source. freqHz should
// be set to LSE frequency.
func Setup(freqHz uint) {
	setup(freqHz)
}

// SetTime sets current calendar time. It does not affect rtos.Nanosec. Only
// time.Now is affected.
func SetTime(t time.Time) {
	up := time.Duration(rtos.Nanosec())
	t = t.Add(-up)
	setStartTime(t)
}

// ISR should be set as irq.RTCAlarm interrupt handler.
func ISR() { isr() }

// Status returns status of RTC. Ok informs whether RTC is configured and works
// properly, set informs whether callendar time was set.
func Status() (ok, set bool) { return status() }
