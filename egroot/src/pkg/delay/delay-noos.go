// +build noos

package delay

import (
	"runtime/noos"
	"syscall"
)

func millisec(ms int) {
	if ms == 0 {
		return
	}
	end := syscall.Uptime() + uint64(ms)*1e6
	te := noos.TickEvent()
	for syscall.Uptime() < end {
		te.Wait()
	}
}
