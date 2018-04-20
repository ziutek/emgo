// +build f030x6 f030x8

package exti

const (
	RTCALR  Lines = 1 << 17 // Real Time Clock Alarm event.
	USB     Lines = 1 << 18 // USB wakeup.
	RTCTTS  Lines = 1 << 19 // RTC Tamper and TimeStamp events.
	RTCWKUP Lines = 1 << 20 // RTC Wakeup event.
)
