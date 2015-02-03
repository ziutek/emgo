// +build noos

package rtos

import "runtime/noos"

func sleepUntil(end uint64) {
	te := noos.TickEvent()
	for Uptime() < end {
		te.Wait()
	}
}