// +build f40_41xxx f429_439xx f746xx

package exti

const (
	RTCALR  Lines = 1 << 17 // Real Time Clock Alarm event.
	OTGFS   Lines = 1 << 18 // USB OTG FS Wakeup event.
	Ether   Lines = 1 << 19 // Ethernet Wakeup event.
	OTGHS   Lines = 1 << 20 // USB OTG HS Wakeup event.
	RTCTTS  Lines = 1 << 21 // RTC Tamper and TimeStamp events.
	RTCWKUP Lines = 1 << 22 // RTC Wakeup event.
	LPTIM1  Lines = 1 << 23 // LPTIM1 asynchronous event.
)
