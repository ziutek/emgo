// Package cmst implements ticking system clock using Cortex-M SysTick timer.
// It is not recomended for low power applications.
package cmst

import (
	"math"
	"nbl"

	"arch/cortexm/scb"
	"arch/cortexm/systick"
)

// Counting ticks is always accurate but requires math.Muldiv. Counting
// nanoseconds requires only ordinary 64-bit multiply and divide but is
// accurate only in some special cases.

var (
	freqHz      uint
	periodTicks uint32
	counter     nbl.Int64
)

func tons(tick int64) int64 {
	return int64(math.Muldiv(uint64(tick), 1e9, uint64(freqHz)))
}

// Setup setups SysTick to work as sytem clock.
//  periodns - number of nanoseconds between ticks (generating PendSV),
//  hz       - frequency of SysTick clock source,
//  external - false: SysTick uses CPU clock; true: SysTick uses external clock.
func Setup(periodns, hz uint, external bool) {
	freqHz = hz
	st := systick.SYSTICK
	if hz == 0 {
		st.CSR.ClearBits(systick.ENABLE | systick.TICKINT)
		return
	}
	periodTicks = uint32((uint64(periodns)*uint64(hz) + 5e8) / 1e9)
	st.RELOAD().Store(systick.RVR_Bits(periodTicks - 1))
	st.CURRENT().Store(0)
	if external {
		st.CLKSOURCE().Clear()
	} else {
		st.CLKSOURCE().Set()
	}
	st.CSR.SetBits(systick.ENABLE | systick.TICKINT)
}

// SetWakeup: see syscall.SetSysClock.
func SetWakeup(t int64) {
}

// Nanosec: see syscall.SetSysClock.
func Nanosec() int64 {
	if freqHz == 0 {
		return 0
	}
	aba := counter.StartLoad()
	for {
		cnt := counter.TryLoad(aba)
		add := periodTicks - uint32(systick.SYSTICK.CURRENT().Load())
		var ok bool
		if aba, ok = counter.CheckLoad(aba); ok {
			return tons(cnt + int64(add))
		}
	}
}

func sysTickHandler() {
	counter.WriterAdd(int64(periodTicks))
	scb.ICSR_Store(scb.PENDSVSET)
}

//c:__attribute__((section(".SysTick")))
var SysTickVector = sysTickHandler
