// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

package syscall

import (
	"builtin"

	"arch/cortexm"
	"arch/cortexm/scb"
)

func schedNext() {
	if cortexm.IPSR()&0xff == 0 {
		// Called from thread mode.
		builtin.Syscall0(SCHEDNEXT)
	} else {
		// Called from ISR.
		scb.ICSR_Store(scb.PENDSVSET)
	}
}
