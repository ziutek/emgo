// +build noos

package delay

import (
	"rtos"
	"runtime/noos"
)

func millisec(ms int) {
	if ms == 0 {
		return
	}
	end := rtos.Uptime() + uint64(ms)*1e6
	te := noos.TickEvent()
	for rtos.Uptime() < end {
		te.Wait()
	}
}
