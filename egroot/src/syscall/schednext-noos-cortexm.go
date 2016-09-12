// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

package syscall

import (
	"internal"

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
			cortexm.SEV()
		} else {
			// Raise PendSV exception.
			SCB.ICSR.Store(scb.PENDSVSET)
		}
	}
}
