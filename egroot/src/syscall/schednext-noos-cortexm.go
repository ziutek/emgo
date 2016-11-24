// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package syscall

import (
	"internal"
	"sync/fence"

	"arch/cortexm"
	"arch/cortexm/scb"
)

func schedNext() {
	switch cortexm.IPSR() & 0xff {
	case 0:
		// Called from thread mode.
		internal.Syscall0(SCHEDNEXT)
	case cortexm.PendSV:
		// Called from PendSV handler when it sends Alarm event.
	default:
		// Called from ISR.
		SCB := scb.SCB
		if SCB.PENDSVACT().Load() != 0 {
			// Wakeup PendSV handler.
			fence.W_SMP() // Complete all writes before wake up other CPUs.
			cortexm.SEV()
		} else {
			// Raise PendSV exception.
			fence.W() // Treat NVIC as external observer of CPU memory write.
			SCB.ICSR.Store(scb.PENDSVSET)
		}
	}
}
