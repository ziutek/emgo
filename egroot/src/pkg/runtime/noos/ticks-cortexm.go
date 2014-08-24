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
	tickPeriod uint32
	ticks2     [2]uint64
	ticksABA   uintptr
)

func setTickPeriod() {
	// Set tick period to 2 ms (500 context switches per second).
	const periodms = 2
	tickPeriod = uint32(sysClk * periodms / 1e3)
	systick.SetReload(tickPeriod - 1)
	systick.Reset()
}

func sysTickHandler() {
	aba := atomic.LoadUintptr(&ticksABA)
	t := ticks2[aba&1]
	aba++
	ticks2[aba&1] = t + 1
	barrier.Memory()
	atomic.StoreUintptr(&ticksABA, aba)
	tickEvent.Send()

	if tasker.onSysTick {
		exce.PendSV.SetPending()
	}
}

func uptime() uint64 {
	var (
		ticks uint64
		cnt   uint32
	)
	aba := atomic.LoadUintptr(&ticksABA)
	for {
		barrier.Compiler()
		cnt = systick.Val()
		ticks = ticks2[aba&1]
		barrier.Compiler()
		aba1 := atomic.LoadUintptr(&ticksABA)
		if aba == aba1 {
			break
		}
		aba = aba1
	}
	return muldiv64(
		ticks*uint64(tickPeriod)+uint64(tickPeriod-cnt),
		1e9, uint64(sysClk),
	)
}

func muldiv64(x, m, d uint64) uint64 {
	divx := x / d
	modx := x - divx*d
	divm := m / d
	modm := m - divm*d
	return divx*m + modx*divm + modx*modm/d
}
