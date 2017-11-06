// +build l476xx

package exti

const (
	OTGFS   Lines = 1 << 17 // USB OTG FS wakeup event.
	RTCALR  Lines = 1 << 18 // Real Time Clock Alarm event.
	RTCTTS  Lines = 1 << 19 // RTC Tamper and TimeStamp events.
	RTCWKUP Lines = 1 << 20 // RTC Wakeup event.
	COMP1   Lines = 1 << 21 // Comparator 1 wakeup event.
	COMP2   Lines = 1 << 22 // Comparator 2 wakeup event.
	I2C1    Lines = 1 << 23 // I2C1 wakeup event.
	I2C2    Lines = 1 << 24 // I2C2 wakeup event.
	I2C3    Lines = 1 << 25 // I2C3 wakeup event.
	USART1  Lines = 1 << 26 // USART1 wakeup event.
	USART2  Lines = 1 << 27 // USART2 wakeup event.
	USART3  Lines = 1 << 28 // USART3 wakeup event.
	USART4  Lines = 1 << 29 // USART4 wakeup event.
	USART5  Lines = 1 << 30 // USART5 wakeup event.
	LPUART1 Lines = 1 << 31 // LPUART1 wakeup event.

	LPTIM1 Lines = 1 << 32 // LPTIM1 event.
	LPTIM2 Lines = 1 << 33 // LPTIM1 event.
	SWPMI1 Lines = 1 << 34 // SWPMI1 wakeup event.
	PVM1   Lines = 1 << 35 // PVM1 wakeup event.
	PVM2   Lines = 1 << 36 // PVM2 wakeup event.
	PVM3   Lines = 1 << 37 // PVM3 wakeup event.
	PVM4   Lines = 1 << 38 // PVM4 wakeup event.
	LCD    Lines = 1 << 39 // LCD wakeup event.
	I2C4   Lines = 1 << 40 // I2C4 wakeup event.
)
