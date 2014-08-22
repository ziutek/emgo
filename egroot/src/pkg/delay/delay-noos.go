// +build noos

package delay

import "runtime/noos"

func millisec(ms int) {
	if ms == 0 {
		return
	}
	end := noos.Uptime() + uint64(ms*1e6)
	te := noos.TickEvent()
	for noos.Uptime() < end {
		te.Wait()
	}
}
