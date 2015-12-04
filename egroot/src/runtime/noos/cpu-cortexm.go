// +build cortexm0 cortexm3 cortexm4

package noos

import (
	"arch/cortexm/scb"
)

func initCPU() {
	// Enable fault handlers.
	(scb.MEMFAULTENA | scb.BUSFAULTENA | scb.USGFAULTENA).Set()
	// Division by zero will cause the UsageFault.
	scb.DIV_0_TRP.Set()
}
