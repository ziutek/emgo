// +build cortexm0 cortexm3 cortexm4 cortexm4f
package noos

import (
	"syscall"

	"arch/cortexm/systick"
)

var tperiodTicks uint32

func setSystimerFreq(hz uint32) {
	if hz == 0 {
		(systick.ENABLE | systick.TICKINT).Clear()
		return
	}
	const periodns = 2e6
	tperiodTicks = uint32(periodns * uint64(hz) / 1e9)
	systick.RELOAD.Store(tperiodTicks - 1)
	systick.CURRENT.Clear()
	(systick.ENABLE | systick.TICKINT | systick.CLKSOURCE).Set()
}

func sysTickHandler() {
	sysclock.WriterAdd(uint64(tperiodTicks))
	syscall.Alarm.Send()
}

func systimer() uint32 {
	return tperiodTicks - systick.CURRENT.Load()
}
