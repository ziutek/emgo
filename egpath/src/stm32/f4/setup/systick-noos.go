// build +noos

package setup

import (
	"cortexm/systick"
	"runtime/noos"
)

func setSystick() {
	if noos.MaxTasks() == 0 {
		return
	}
	// Set tick period to 2 ms (500 context switches per second).
	const period = 2
	systick.SetReload(uint32(SysClk*period/1e3 - 1))
	noos.SetTickPeriod(period)
}
