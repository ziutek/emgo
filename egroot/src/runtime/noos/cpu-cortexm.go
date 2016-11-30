// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package noos

import (
	"internal"
	"unsafe"

	"arch/cortexm"
	"arch/cortexm/acc"
	"arch/cortexm/fpu"
	"arch/cortexm/scb"
)

func initCPU() {
	SCB := scb.SCB
	FPU := fpu.FPU
	ACC := acc.ACC
	// Enable fault handlers.
	SCB.SHCSR.SetBits(scb.MEMFAULTENA | scb.BUSFAULTENA | scb.USGFAULTENA)
	// Division by zero will cause the UsageFault.
	SCB.DIV_0_TRP().Set()
	if hasFPU {
		// Enable FPU.
		FPU.CP10().Store(fpu.CPACFULL << fpu.CP10n)
		FPU.FPCCR.Store(fpu.ASPEN | fpu.LSPEN)
	}
	// Move exception vectors to ITCM RAM if available.
	if vtor := SCB.VTOR.Load(); vtor != 0 {
		cr := ACC.ITCMCR.Load()
		if cr&acc.ITCMSZ != 0 {
			if cr&acc.ITCMEN == 0 {
				ACC.ITCMCR.Store(cr | acc.ITCMEN)
				cortexm.DSB()
			}
			internal.Memmove(nil, unsafe.Pointer(uintptr(vtor)), vectorsSize())
			SCB.VTOR.Store(0)
			cortexm.DSB()
		}
	}
}
