// build +noos

package setup

import (
	"runtime/noos"
	"syscall"
)

func sysClkChanged() {
	if noos.MaxTasks() == 0 {
		return
	}
	syscall.SetSysClock(SysClk)
}
