// Package cmst implements system clock using Cortex-M SysTick timer.
package cmst

import (
	"math"
	"nbl"

	"arch/cortexm/scb"
	"arch/cortexm/systick"
)

var (
	freqHz      uint   // Hz.
	periodTicks uint32 // Ticks.
	counter     nbl.Int64
	wakeup      nbl.Int64
)

func Setup(periodns, hz uint, external bool) {
	freqHz = hz
	if hz == 0 {
		(systick.ENABLE | systick.TICKINT).Clear()
		return
	}
	periodTicks = uint32((uint64(periodns)*uint64(hz) + 5e8) / 1e9)
	systick.RELOAD.Store(periodTicks - 1)
	systick.CURRENT.Clear()
	cfg := systick.ENABLE | systick.TICKINT
	if !external {
		cfg |= systick.CLKSOURCE
	}
	cfg.Set()
}

// SetWakeup: see syscall.SetSysClock.
func SetWakeup(t int64) {
	wakeup.WriterStore(t)
}

// Uptime: see syscall.SetSysClock.
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
