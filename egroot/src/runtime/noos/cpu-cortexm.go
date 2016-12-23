// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package noos

import (
	"internal"
	"unsafe"

	"arch/cortexm"
	"arch/cortexm/acc"
	//"arch/cortexm/cmt"
	"arch/cortexm/fpu"
	"arch/cortexm/pft"
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
	ACC := acc.ACC
	if useITCM {
		// Move exception vectors to ITCM RAM if available.
		if vtor := SCB.VTOR.Load(); vtor != 0 {
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
		//CMT := cmt.CMT
		PFT := pft.PFT
		PFT.CSSELR.Store(0) // Select data cache size.
		cortexm.DSB()
		ccsidr := PFT.CCSIDR.Load()
		lsiz := 4 << uint(ccsidr.Field(pft.LineSize))
		nways := ccsidr.Field(pft.Associativity) - 1
		nsets := ccsidr.Field(pft.NumSets) - 1
		_ = lsiz
		_ = nways
		_ = nsets
		// Enable cache in write-through mode on shareable regions.
		ACC.CACR.Store(acc.SIWT)
	}
}
