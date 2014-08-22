// build +noos

package setup

import "runtime/noos"

func sysClkChanged() {
	if noos.MaxTasks() == 0 {
		return
	}
	noos.SetSysClk(SysClk)
}
