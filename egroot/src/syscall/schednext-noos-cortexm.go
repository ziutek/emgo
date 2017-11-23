// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package syscall

import (
	"sync/fence"

	"arch/cortexm"
	"arch/cortexm/scb"
)

func schedNext() {
	switch cortexm.IPSR() & 0xff {
	case 0:
		// Called from thread mode.
		SchedYield()
	case cortexm.PendSV: // Called from PendSV
		for {
			// This should not happen!
		}
	default: // Called from ISR
		// Raise PendSV exception.
		fence.W()     // Treat NVIC as external observer of CPU memory write.
		cortexm.SEV() // See ARM Errata 563915 or STM32F10xx Errata 1.1.2.
		scb.SCB.ICSR.Store(scb.PENDSVSET)
	}
}
