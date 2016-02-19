// +build cortexm4f

package noos

import (
	"arch/cortexm/fpu"
	"arch/cortexm/scb"
)

func initCPU() {
	SCB := scb.SCB
	FPU := fpu.FPU
	// Enable fault handlers.
	SCB.SHCSR.SetBits(scb.MEMFAULTENA | scb.BUSFAULTENA | scb.USGFAULTENA)
	// Division by zero will cause the UsageFault.
	SCB.DIV_0_TRP().Set()
	// Enable FPU.
	FPU.CP10().Store(fpu.CPACFULL << fpu.CP10n)
	FPU.FPCCR.Store(fpu.ASPEN | fpu.LSPEN)
}
