// +build noos

package setup

import (
	"runtime/noos/sysclk/cmst"
	"syscall"
)

// UseSysTick setups and uses Cortex-M SYSTICK timer as system clock.
func UseSysTick() {
	lev, _ := syscall.SetPrivLevel(0)
	cmst.Setup(2e6, AHBClk/8, true)
	syscall.SetPrivLevel(lev)
	syscall.SetSysClock(cmst.Nanosec, cmst.SetWakeup)
}

// UseRTC setups and uses STM32 real time clock as system clock. It returns true if
// clock
func UseRTC(freq uint) {
	useRTC(freq)
}
