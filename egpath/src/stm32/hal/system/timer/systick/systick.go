// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package systick

import (
	"runtime/noos/timer/cmst"
	"syscall"

	"stm32/hal/system"
)

// Setup setups and uses Cortex-M SYSTICK as system timer.
func Setup() {
	lev, _ := syscall.SetPrivLevel(0)
	cmst.Setup(2e6, system.AHB.Clock()/8, true)
	syscall.SetPrivLevel(lev)
	syscall.SetSysTimer(cmst.Nanosec, cmst.SetWakeup)
}
