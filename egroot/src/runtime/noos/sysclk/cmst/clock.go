// Package cmst implements ticking system clock using Cortex-M SysTick timer.
// It is not recomended for low power applications.
package cmst

import (
	"math"
	"nbl"

	"arch/cortexm/scb"
	"arch/cortexm/systick"
)

var (
	freqHz      uint
	periodTicks uint32
	counter     nbl.Int64
)

// Setup setups SysTick to work as sytem clock.
//  periodns - number of nanoseconds between ticks (generating PendSV),
//  hz       - frequency of SysTick clock source,
//  external - false: SysTick uses CPU clock; true: SysTick uses external clock.
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

func sysTickHandler() {
	counter.WriterAdd(int64(periodTicks))
	scb.ICSR_Store(scb.PENDSVSET)
}

//c:__attribute__((section(".SysTick")))
var SysTickVector = sysTickHandler
