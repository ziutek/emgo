// Package cmst implements system clock using Cortex-M SysTick timer.
package cmst

import (
	"math"
	"nbl"

	"arch/cortexm/scb"
	"arch/cortexm/systick"
)

var (
	freqHz      uint32 // Hz.
	periodTicks uint32 // Ticks.
	counter     nbl.Int64
	wakeup      nbl.Int64
)

// SetFreq setups and starts/stops the system clock. hz == 0 means stop the
// clock. hz > 0 informs system clock about freqency of SysTick clock source
// and starts system clock. If external is true SysTick is configured to use
// external source, instead it uses CPU clock.
func SetFreq(hz uint32, external bool) {
	freqHz = hz
	if hz == 0 {
		(systick.ENABLE | systick.TICKINT).Clear()
		return
	}
	const periodns = 2e6
	periodTicks = uint32(periodns * uint64(hz) / 1e9)
	systick.RELOAD.Store(periodTicks - 1)
	systick.CURRENT.Clear()
	cfg := systick.ENABLE | systick.TICKINT
	if !external {
		cfg |= systick.CLKSOURCE
	}
	cfg.Set()
}

// SetWakeup asks timer to run runtime scheduler at uptime t nanoseconds using
// PendSV exception. t is only a hint and scheduler can be run at any uptime
// less or equal to t or even greather than t, if t is less than current uptime.
// Notice that scheduler can be also run at any time by other PendSV source.
func SetWakeup(t int64) {
	wakeup.WriterStore(t)
}

// Uptime returns the time elapsed since the start of the system clock.
func Uptime() int64 {
	if freqHz == 0 {
		return 0
	}
	aba := counter.StartLoad()
	for {
		cnt := uint64(counter.TryLoad(aba))
		add := periodTicks - systick.CURRENT.Load()
		var ok bool
		if aba, ok = counter.CheckLoad(aba); ok {
			return int64(math.Muldiv(cnt+uint64(add), 1e9, uint64(freqHz)))
		}
	}
}

func SysTickHandler() {
	counter.WriterAdd(int64(periodTicks))
	scb.ICSR_Store(scb.PENDSVSET)
}

//c:__attribute__((section(".SysTick")))
var SysTickVector = SysTickHandler
