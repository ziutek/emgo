// +build cortexm0 cortexm3 cortexm4

package noos

import (
	"arch/cortexm/scb"
)

func initCPU() {
	// Division by zero and unaligned access will cause the UsageFault.
	scb.CCR.SetBits(scb.DIV_0_TRP | scb.UNALIGN_TRP)
}
