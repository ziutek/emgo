// +build cortexm4f

package noos

import (
	"arch/cortexm/scb"
	"arch/cortexm/fpu"
)

func initCPU() {
	scb.R.CCR.SetBit(scb.DIV_0_TRP) // Enable div by 0 exception.
	fpu.SetAccess(fpu.Full)
	fpu.SetSP(fpu.AutoSP | fpu.LazySP)
}