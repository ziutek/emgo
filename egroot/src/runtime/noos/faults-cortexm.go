// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm"
)

func NMIHandler() {
	cortexm.BKPT(0)
}

func FaultHandler()
