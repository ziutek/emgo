// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package exti

const (
	RTCALR Lines = 1 << 17 // Real Time Clock Alarm event.
	USB    Lines = 1 << 18 // USB wakeup.
	Ether  Lines = 1 << 19 // Ethernet wakeup.
)
