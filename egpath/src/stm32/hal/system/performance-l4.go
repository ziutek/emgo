// +build l476xx

package system

import (
	"stm32/hal/raw/rcc"
)

func Setup() {
	RCC := rcc.RCC

	// Reset RCC clock configuration.
	RCC.MSION().Set()
	for RCC.MSIRDY().Load() == 0 {
		// Wait for MSI...
	}
	RCC.CFGR.Store(0)
	RCC.CR.ClearBits(rcc.HSION | rcc.HSEON | rcc.CSSON | rcc.PLLON | rcc.HSEBYP)
	RCC.CIER.Store(0) // Disable clock interrupts.

	clock[Core] = 4e6
	clock[AHB] = 4e6
	clock[APB1] = 4e6
	clock[APB2] = 4e6
}
