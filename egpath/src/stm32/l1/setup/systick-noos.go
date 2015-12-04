// build +noos

package setup

import (
	"runtime/noos"
	"syscall"
)

func sysclkChanged() {
	if noos.MaxTasks() == 0 {
		return
	}
	syscall.SetSysClock(SysClk, nil)
}
