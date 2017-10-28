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
	default: // Called from ISR (raise PendSV exception).
		fence.W() // Treat NVIC as external observer of CPU memory write.
		scb.SCB.ICSR.Store(scb.PENDSVSET)
	}
}
