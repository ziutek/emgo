// +build noos

package delay

import (
	"log"
	"runtime/noos"
)

func millisec(ms int) {
	if noos.MaxTasks() == 0 {
		log.Panic("no support for delay.Millisec")
	}
	period := noos.TickPeriod()
	if period == 0 {
		log.Panic("tick period not configured in runtime/noos")
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
