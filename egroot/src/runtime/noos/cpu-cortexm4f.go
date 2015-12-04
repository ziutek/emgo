// +build cortexm4f

package noos

import (
	"arch/cortexm/fpu"
	"arch/cortexm/scb"
)

func initCPU() {
	// Enable fault handlers.
	(scb.MEMFAULTENA | scb.BUSFAULTENA | scb.USGFAULTENA).Set()
	// Division by zero will cause the UsageFault.
	scb.DIV_0_TRP.Set()
	// Enable FPU.
	scb.CP10.Store(scb.AccessFull)
	fpu.FPCCR_Store(fpu.ASPEN | fpu.LSPEN)
}
