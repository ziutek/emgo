// Package fmc.
package fmc

import (
	"stm32/hal/raw/rcc"
)

// EnableClock enables clock for FMC/FSMC.
func EnableClock(lp bool) {
	RCC := rcc.RCC
	if lp {
		RCC.FMCLPEN().Set()
	} else {
		RCC.FMCLPEN().Clear()

	}
	RCC.FMCEN().Set()
}

// DisableClock disables clock for FMC/FSMC.
func DisableClock() {
	rcc.RCC.FMCEN().Clear()
}
