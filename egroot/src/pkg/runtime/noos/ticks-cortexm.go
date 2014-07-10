// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"cortexm/exce"
	"cortexm/systick"
	"sync/atomic"
	"sync/barrier"
)

func sysTickStart() {
	// Defaults:
	// One context switch per 1e6 SysTicks (70/s for 70 Mhz, 168/s for 168 MHz)
	systick.SetReload(1e6 - 1)
	systick.WriteFlags(systick.Enable | systick.TickInt | systick.ClkCPU)
}

var (
	ticks    [2]uint64
	ticksABA uintptr
)

func sysTickHandler() {
	aba := atomic.LoadUintptr(&ticksABA)
	t := ticks[aba&1]
	aba++
	ticks[aba&1] = t + 1
	barrier.Memory()
	atomic.StoreUintptr(&ticksABA, aba)
	tickEvent.Send()

	if tasker.onSysTick {
		exce.PendSV.SetPending()
	}
}

func loadTicks() uint64 {
	aba := atomic.LoadUintptr(&ticksABA)
	for {
		barrier.Compiler()
		t := ticks[aba&1]
		barrier.Compiler()
		aba1 := atomic.LoadUintptr(&ticksABA)
		if aba == aba1 {
			return t
		}
		aba = aba1
	}
}
