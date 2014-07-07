// +build noos

package delay

import "runtime/noos"

func millisec(ms int) {
	if noos.MaxTasks() == 0 {
		panic("no support for delay.Millisec")
	}
	period := noos.TickPeriod()
	if period == 0 {
		panic("tick period not configured in runtime/noos")
	}
	dt := ms / period
	if dt == 0 {
		return
	}
	to := noos.Ticks() + uint64(dt)
	for {
		noos.TickEvent().Wait()
		if noos.Ticks() >= to {
			return
		}
	}
}
