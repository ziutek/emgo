// +build cortexm0 cortexm3 cortexm4

package noos

import (
	"arch/cortexm/scb"
)

func initCPU() {
	scb.R.CCR.SetBit(scb.DIV_0_TRP) // Enable div by 0 exception.
}