// +build cortexm3 cortexm4 cortexm4f

package systick

import (
	"runtime/noos/clk/cmst"
	"syscall"

	"stm32/hal/system"
)

// Setup setups and uses Cortex-M SYSTICK timer as OS clock.
func Setup() {
	lev, _ := syscall.SetPrivLevel(0)
	cmst.Setup(2e6, system.AHBClk/8, true)
	syscall.SetPrivLevel(lev)
	syscall.SetSysClock(cmst.Nanosec, cmst.SetWakeup)
}
