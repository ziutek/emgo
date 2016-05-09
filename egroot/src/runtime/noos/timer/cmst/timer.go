// Package cmst implements ticking OS timer using Cortex-M SysTick timer.
// It is not recomended for low power applications.
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

var (
	freqHz      uint
	periodTicks uint32
	counter     nbl.Int64
)

func tons(tick int64) int64 {
	return int64(math.Muldiv(uint64(tick), 1e9, uint64(freqHz)))
}

// Setup setups SysTick to work as sytem timer.
//  periodns - number of nanoseconds between ticks (generating PendSV),
//  hz       - frequency of SysTick clock source,
//  external - false: SysTick uses CPU clock; true: SysTick uses external clock.
// Setup must be run in privileged mode.
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
	// Set priority for SysTick exception higher SVCall proprity, to ensure
	// that any user of rtos.Nanosec observes both SYSTICK.CURRENT reset and
	// counter increment as one atomic operation.
	spnum := cortexm.PrioStep * cortexm.PrioNum
	prio := cortexm.PrioLowest + spnum*3/4
	scb.SCB.PRI_SysTick().Store(scb.PRI_SysTick.J(prio))
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
	aba := counter.ABA()
	for {
		cnt := counter.TryLoad(aba)
		add := periodTicks - uint32(systick.SYSTICK.CURRENT().Load())
		var ok bool
		if aba, ok = counter.CheckABA(aba); ok {
			return tons(cnt + int64(add))
		}
	}
}

func sysTickHandler() {
	counter.WriterAdd(int64(periodTicks))
	syscall.SchedNext()
}

//emgo:const
//c:__attribute__((section(".SysTick")))
var SysTickVector = sysTickHandler
