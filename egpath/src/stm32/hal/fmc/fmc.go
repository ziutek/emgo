// Package fmc.
package fmc

import (
	"stm32/hal/raw/rcc"
)

// EnableClock enables clock for FMC/FSMC.
func EnableClock(lp bool) {
	RCC := rcc.RCC
	if lp {
		RCC.FMCLPEN().AtomicSet()
	} else {
		RCC.FMCLPEN().AtomicClear()

	}
	RCC.FMCEN().AtomicSet()
}

// DisableClock disables clock for FMC/FSMC.
func DisableClock() {
	rcc.RCC.FMCEN().AtomicClear()
}
