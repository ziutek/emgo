// +build cortexm4f

package noos

import (
	"arch/cortexm/scb"
	"arch/cortexm/fpu"
)

func initCPU() {
	// Division by zero and unaligned access will cause the UsageFault.
	scb.CCR.SetBits(scb.DIV_0_TRP | scb.UNALIGN_TRP)
	
	fpu.SetAccess(fpu.Full)
	fpu.SetSP(fpu.AutoSP | fpu.LazySP)
}