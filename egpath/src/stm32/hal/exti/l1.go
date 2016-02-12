// +build l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package exti

const (
	USBFS   Lines = 1 << 18 // USB Device FS wakeup event.
	RTCTTS  Lines = 1 << 19 // RTC Tamper and TimeStamp events.
	RTCWKUP Lines = 1 << 20 // RTC Wakeup event.
	COMP1         = 1 << 21 // Comparator 1 wakeup event.
	COMP2         = 1 << 22 // Comparator 2 wakeup event.
	COMPCA        = 1 << 23 // Channel acquisition interrupt.
)
