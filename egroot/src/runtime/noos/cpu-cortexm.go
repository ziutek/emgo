// +build cortexm0 cortexm3 cortexm4

package noos

import (
	"arch/cortexm/scb"
)

func initCPU() {
	SCB := scb.SCB
	// Enable fault handlers.
	SCB.SHCSR.SetBits(scb.MEMFAULTENA | scb.BUSFAULTENA | scb.USGFAULTENA)
	// Division by zero will cause the UsageFault.
	SCB.DIV_0_TRP().Set()
}
