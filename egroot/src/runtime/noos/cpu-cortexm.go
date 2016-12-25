// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package noos

import (
	"bits"
	"internal"
	"unsafe"

	"arch/cortexm"
	"arch/cortexm/acc"
	"arch/cortexm/cmt"
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
		CMT := cmt.CMT
		PFT := pft.PFT

		// Invalidate data cache.
		PFT.CSSELR.Store(0) // Select L1 cache size info.
		cortexm.DSB()
		csi := PFT.CCSIDR.Load() // Load L1 cache size info.

		maxset := uint32(csi.Field(pft.NumSets))       // Max. set number.
		maxway := uint32(csi.Field(pft.Associativity)) // Max. way number.
		log2bpl := uint(csi.Field(pft.LineSize) + 4)   // Log2(bytes per line).
		wayshift := bits.LeadingZeros32(uint32(maxway))

		for set := uint32(0); set <= maxset; set++ {
			for way := uint32(0); way <= maxway; way++ {
				CMT.DCISW.U32.Store(way<<wayshift | set<<log2bpl)
			}
		}

		// Invalidate instruction cache.
		CMT.ICIALLU.Store(0)

		// Use cache in write-through mode on shareable regions.
		ACC.CACR.Store(acc.SIWT)

		cortexm.DSB()
		cortexm.ISB()

		// Enable data and instruction cache.
		SCB.CCR.SetBits(scb.DC | scb.IC)

		cortexm.DSB()
		cortexm.ISB()
	}
}

/*lsiz := 4 << uint(ccsidr.Field(pft.LineSize))
nways := ccsidr.Field(pft.Associativity) - 1
nsets := ccsidr.Field(pft.NumSets) - 1
_ = lsiz
_ = nways
_ = nsets*/
