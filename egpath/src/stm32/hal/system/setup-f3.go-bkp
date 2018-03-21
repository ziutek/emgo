// +build f303xe

package system

import (
	"stm32/hal/raw/flash"
	"stm32/hal/raw/rcc"
)

// SetupPLL setups MCU for best performance (prefetch on, minimum allowed Flash
// latency) using integrated PLL as system clock source.
//
// Osc is freqency of external resonator in MHz. Allowed values: 4 to 32 MHz.
// Use 0 to select internal HSI oscilator (8 MHz / 2) as system clock source.
//
// Div and mul determine the system clock frequency according to the formula:
//
// When osc == 0 or:
//
//  SysClk = 8e6 / div * mul [Hz]
//
// div can be 1..16 (303xD, 303xE, 398xE) or must be 2 (other).
//
// When osc != 0.
//
//  SysClk = osc * 1e6 / div * mul [Hz]
//
// div can be 1..16 (303xD, 303xE, 398xE) or 2..16 (other).
//
// Mul can be 2..16.
//
// USB requires HSE as PLL clock source and SysClk set to 48e6 or 72e6 Hz.
func SetupPLL(osc, div, mul int) {
	RCC := rcc.RCC

	// Reset RCC clock configuration.
	RCC.HSION().Set()
	for RCC.HSIRDY().Load() == 0 {
		// Wait for HSI...
	}
	RCC.CFGR.Store(0)
	for RCC.SWS().Load() != rcc.SWS_HSI {
		// Wait for system clock setup...
	}
	RCC.CR.ClearBits(rcc.HSEON | rcc.CSSON | rcc.PLLON | rcc.HSEBYP)
	RCC.CIR.Store(0) // Disable clock interrupts.

	// Calculate system clock.
	if osc != 0 && (osc < 4 || osc > 32) {
		panic("bad HSE osc freq")
	}
	if mul < 2 || mul > 16 {
		panic("bad PLL multipler")
	}
	var sysclk uint
	if osc == 0 {
		sysclk = HSIClk / uint(div) * uint(mul) // Hz
	} else {
		// HSE needs milliseconds to stabilize, so enable it now.
		RCC.HSEON().Set()
		sysclk = uint(osc) * 1e6 / uint(div) * uint(mul) // Hz
	}
	ahbclk := sysclk
	cfgr := rcc.CFGR(mul-2) << rcc.PLLMULn
	var apb1clk uint
	switch {
	case ahbclk <= 1*maxAPB1Clk:
		cfgr |= rcc.PPRE1_DIV1
		apb1clk = ahbclk / 1
	case ahbclk <= 2*maxAPB1Clk:
		cfgr |= rcc.PPRE1_DIV2
		apb1clk = ahbclk / 2
	case ahbclk <= 4*maxAPB1Clk:
		cfgr |= rcc.PPRE1_DIV4
		apb1clk = ahbclk / 4
	case ahbclk <= 8*maxAPB1Clk:
		cfgr |= rcc.PPRE1_DIV8
		apb1clk = ahbclk / 8
	default:
		cfgr |= rcc.PPRE1_DIV16
		apb1clk = ahbclk / 16
	}
	clock[Core] = sysclk
	clock[AHB] = ahbclk
	clock[APB1] = apb1clk
	clock[APB2] = ahbclk
	if sysclk <= 48e6 {
		cfgr |= rcc.USBPRE
	}
	// Setup Flash.
	FLASH := flash.FLASH
	latency := flash.ACR((sysclk-1)/24e6) << flash.LATENCYn
	FLASH.ACR.SetBits(flash.PRFTBE | latency)
	// Setup PLL.
	if osc == 0 {
		// Div == 2 for HSI can be selected in compatible way: PLLSRC = 0.
		if div != 2 {
			cfgr |= rcc.PLLSRC_HSI_PREDIV // PLLSRC = 1
		}
	} else {
		cfgr |= rcc.PLLSRC_HSE_PREDIV // PLLSRC = 2
		for RCC.HSERDY().Load() == 0 {
			// Wait for HSE...
		}
	}
	RCC.CFGR.Store(cfgr)
	RCC.CFGR2.Store(rcc.CFGR2(div-1) << rcc.PREDIVn) // Must be after CFGR.
	RCC.PLLON().Set()
	for RCC.PLLRDY().Load() == 0 {
		// Wait for PLL...
	}
	// Change system clock source to PLL.
	RCC.SW().Store(rcc.SW_PLL)
	for RCC.SWS().Load() != rcc.SWS_PLL {
		// Wait for system clock setup...
	}
	if osc != 0 {
		RCC.HSION().Clear()
	}
}
