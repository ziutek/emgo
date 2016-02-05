// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

package syscall

import (
	"builtin"

	"arch/cortexm"
	"arch/cortexm/scb"
)

func schedNext() {
	switch cortexm.IPSR() & 0xff {
	case 0:
		// Called from thread mode.
		builtin.Syscall0(SCHEDNEXT)
	case cortexm.PendSV:
		// Called from PendSV handler when it sends Alarm event.
	default:
		// Called from ISR.
		if scb.PENDSVACT.Load() != 0 {
			// Wakeup PendSV handler.
			cortexm.SEV()
		} else {
			// Rise PendSV exception.
			scb.ICSR_Store(scb.PENDSVSET)
		}
	}
}
