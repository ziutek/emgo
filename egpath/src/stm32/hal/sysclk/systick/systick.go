// +build cortexm3 cortexm4 cortexm4f

package systick

import (
	"runtime/noos/sysclk/cmst"
	"syscall"

	"stm32/hal/setup"
)

// UseSysTick setups and uses Cortex-M SYSTICK timer as system clock.
func Setup() {
	lev, _ := syscall.SetPrivLevel(0)
	cmst.Setup(2e6, setup.AHBClk/8, true)
	syscall.SetPrivLevel(lev)
	syscall.SetSysClock(cmst.Nanosec, cmst.SetWakeup)
}
