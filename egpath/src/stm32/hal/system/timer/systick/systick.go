// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package systick

import (
	"runtime/noos/timer/cmst"
	"syscall"

	"stm32/hal/system"
)

// Setup setups Cortex-M SYSTICK as system timer. This is ticking timer. Use
// tickless timer if available. SYSTICK runs scheduler every periodns
// nanoseconds.
func Setup(periodns uint) {
	lev, _ := syscall.SetPrivLevel(0)
	cmst.Setup(periodns, system.AHB.Clock()/8, true)
	syscall.SetPrivLevel(lev)
	syscall.SetSysTimer(cmst.Nanosec, cmst.SetWakeup)
}
