// +build cortexm4f

package noos

import "arch/cortexm/fpu"

func initCPU() {
	fpu.SetAccess(fpu.Full)
	fpu.SetSP(fpu.AutoSP | fpu.LazySP)
}