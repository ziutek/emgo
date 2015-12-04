package noos

import (
	"math"
	"nbl"
)

var (
	sysclock   nbl.Uint64
	sysrtc     func() uint32
	sysclockHz uint32
)

func uptime() uint64 {
	aba := sysclock.StartLoad()
	for {
		cnt := sysclock.TryLoad(aba)
		rtc := sysrtc()
		var ok bool
		if aba, ok = sysclock.CheckLoad(aba); ok {
			return math.Muldiv(cnt+uint64(rtc), 1e9, uint64(sysclockHz))
		}
	}

}
