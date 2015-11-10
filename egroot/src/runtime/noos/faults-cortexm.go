// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm"
)

//emgo:noinline
func nmiHandler() {
	cortexm.BKPT(0)
}

func FaultHandler()
