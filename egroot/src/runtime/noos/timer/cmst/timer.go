// Package cmst implements ticking OS timer using Cortex-M SysTick timer. It can
// be used only in uniprocessor system and is not recomended for low power
// applications.
package cmst

import (
	"math"
	"nbl"
	"syscall"

	"arch/cortexm"
	"arch/cortexm/scb"
	"arch/cortexm/systick"
)

// Counting ticks is always accurate but requires math.Muldiv. Counting
// nanoseconds requires only ordinary 64-bit multiply and divide but is
// accurate only in some special cases.

type globals struct {
	counter     nbl.Int64
	freqHz      uint
	periodTicks uint32
}

var g globals

func ticktons(tick int64) int64 {
	return int64(math.MulDiv(uint64(tick), 1e9, uint64(g.freqHz)))
}

// Setup setups SysTick to work as sytem timer.
//  periodns - number of nanoseconds between ticks,
//  hz       - frequency of SysTick clock source,
//  external - false: SysTick uses CPU clock; true: SysTick uses external clock.
// Setup must be run in privileged mode.
func Setup(periodns, hz uint, external bool) {
	g.freqHz = hz
	st := systick.SYSTICK
	if hz == 0 {
		st.CSR.ClearBits(systick.ENABLE | systick.TICKINT)
		return
	}
	g.periodTicks = uint32((uint64(periodns)*uint64(hz) + 5e8) / 1e9)
	st.RELOAD().Store(systick.RVR_Bits(g.periodTicks - 1))
	st.CURRENT().Store(0)
	var clksrc systick.CSR_Bits
	if !external {
		clksrc = systick.CLKSOURCE
	}
	// Set priority for SysTick exception higher SVCall proprity, to ensure
	// that any user of rtos.Nanosec observes both SYSTICK.CURRENT reset and
	// counter increment as one atomic operation.
	spnum := cortexm.PrioStep * cortexm.PrioNum
	prio := cortexm.PrioLowest + spnum*3/4
	scb.SCB.PRI_SysTick().Store(scb.PRI_SysTick.J(prio))
	st.CSR.Store(systick.ENABLE | systick.TICKINT | clksrc)
}

// SetWakeup: see syscall.SetSysTimer.
func SetWakeup(ns int64) {}

// Nanosec: see syscall.SetSysClock.
func Nanosec() int64 {
	if g.freqHz == 0 {
		return 0
	}
	aba := g.counter.ABA()
	for {
		cnt := g.counter.TryLoad(aba)
		add := g.periodTicks - uint32(systick.SYSTICK.CURRENT().Load())
		var ok bool
		if aba, ok = g.counter.CheckABA(aba); ok {
			return ticktons(cnt + int64(add))
		}
	}
}

func sysTickHandler() {
	g.counter.WriterAdd(int64(g.periodTicks))
	syscall.SchedNext()
}

//emgo:const
//c:__attribute__((section(".SysTick")))
var SysTickVector = sysTickHandler
