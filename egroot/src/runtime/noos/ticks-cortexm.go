// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm/exce"
	"arch/cortexm/systick"
	"sync/atomic"
	"sync/fence"
)

var (
	ticks2     [2]uint64
	ticksABA   uintptr
	sysClk     uint32
	tickPeriod uint32
)

func sysTickStart() {
	// Defaults:
	// One context switch per 1e6 SysTicks (70/s for 70 Mhz, 168/s for 168 MHz)
	systick.SetReload(1e6 - 1)
	systick.Reset()
	systick.StoreFlags(systick.Enable | systick.TickInt | systick.ClkCPU)
}

func setTickPeriod() {
	// Set tick period to 2 ms (500 context switches per second).
	const periodms = 2
	tickPeriod = sysClk * periodms / 1e3
	systick.SetReload(tickPeriod - 1)
	systick.Reset()
}

func sysTickHandler() {
	updateTicks2(tickPeriod)
	tickEvent.Send()
	exce.PendSV.SetPending()
}

func updateTicks2(delta uint32) {
	aba := atomic.LoadUintptr(&ticksABA)
	t := ticks2[aba&1]
	aba++
	ticks2[aba&1] = t + uint64(delta)
	fence.Memory()
	atomic.StoreUintptr(&ticksABA, aba)
}

func uptime() uint64 {
	var (
		ticks uint64
		cnt   uint32
	)
	aba := atomic.LoadUintptr(&ticksABA)
	for {
		fence.Compiler()
		cnt = systick.Val()
		ticks = ticks2[aba&1]
		fence.Compiler()
		aba1 := atomic.LoadUintptr(&ticksABA)
		if aba == aba1 {
			break
		}
		aba = aba1
	}
	return muldiv64(
		ticks+uint64(tickPeriod-cnt),
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
