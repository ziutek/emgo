// Package rtc implements tickless system timer using real time counter.
package rtc

import(
	"nrf51/hal/rtc"
)

var sysRTC *rtc.Periph

// Setup setups st as system timer. 
func Setup(st *rtc.Periph) {
	sysRTC = st
	st.TASK(rtc.STOP).Trigger()
	st.TASK(rtc.CLEAR).Trigger()
	st.SetPRESCALER(0) // 0 means 512 s period.
	
}