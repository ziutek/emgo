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
	// Enable fault handlers.
	SCB.SHCSR.SetBits(scb.MEMFAULTENA | scb.BUSFAULTENA | scb.USGFAULTENA)
	// Division by zero will cause the UsageFault.
	SCB.DIV_0_TRP().Set()
	if useFPU {
		FPU := fpu.FPU
		FPU.CP10().Store(fpu.CPACFULL << fpu.CP10n)
		FPU.FPCCR.Store(fpu.ASPEN | fpu.LSPEN)
	}
	if useITCM {
		// Move exception vectors to ITCM RAM if available.
		if vtor := SCB.VTOR.Load(); vtor != 0 {
			ACC := acc.ACC
			cr := ACC.ITCMCR.Load()
			if cr&acc.ITCMSZ != 0 {
				if cr&acc.ITCMEN == 0 {
					ACC.ITCMCR.Store(cr | acc.ITCMEN)
					cortexm.DSB()
				}
				dst := unsafe.Pointer(uintptr(0))
				src := unsafe.Pointer(uintptr(vtor))
				internal.Memmove(dst, src, vectorsSize())
				SCB.VTOR.Store(0)
				cortexm.DSB()
			}
		}
	}
	if useL1Cache {
	
	}
}
