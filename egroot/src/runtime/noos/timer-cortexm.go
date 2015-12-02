// +build cortexm0 cortexm3 cortexm4 cortexm4f
package noos

import (
	"syscall"

	"arch/cortexm/systick"
)

var timerPeriod_cnt uint32

func setTimerFreq(hz uint32) {
	if hz == 0 {
		(systick.ENABLE | systick.TICKINT).Clear()
		return
	}
	const timerPeriod_ns = 2e6
	timerPeriod_cnt = uint32(timerPeriod_ns * uint64(hz) / 1e9)
	systick.RELOAD.Store(timerPeriod_cnt - 1)
	systick.CURRENT.Clear()
	(systick.ENABLE | systick.TICKINT | systick.CLKSOURCE).Set()
}

func sysTickHandler() {
	sysCounter.Add(uint64(timerPeriod_cnt))
	syscall.Alarm.Send()
}

func timerCnt() uint64 {
	return uint64(timerPeriod_cnt - systick.CURRENT.Load())
}
