// +build l1xx_md l1xx_mdp l1xx_hd l1xx_xl

// Package setup allows to easy setup MCU for typical use.
//
// Clock setup
//
// Goal is to provide 48 MHz for USB so PLLCLK must be set to 96 MHz because
// USBCLK = PLLCLK / 2.
//
// System Clock is derived from PLLCK as follows:
//
//  SYSCLK = PLLCK / PLLDIV
//
package system

import (
	"stm32/hal/raw/flash"
	"stm32/hal/raw/rcc"
)

// SetupPLL setups MCU for best performance (prefetch on, 64-bit flash access)
// using integrated PLL as system clock source.
//
// Osc is freqency of external resonator in MHz. Allowed values: 2, 3, 4, 6, 8,
// 12, 16, 24. Use 0 to select internal HSI oscilator as system clock source.
//
// sdiv is system clock divider. Allowed values: 2, 3, 4. sdiv determine the
// system clock frequency according to the formula:
//
//  SysClk = 96e6 / sdiv [Hz]
//
func SetupPLL(osc, sdiv int) {
	RCC := rcc.RCC

	// Reset RCC clock configuration.
	RCC.MSION().Set()
	for RCC.MSIRDY().Load() == 0 {
		// Wait for MSI...
	}
	RCC.CFGR.Store(0)
	RCC.CR.ClearBits(rcc.HSION | rcc.HSEON | rcc.CSSON | rcc.PLLON | rcc.HSEBYP)
	RCC.CIR.Store(0) // Disable clock interrupts.

	switch sdiv {
	case 2, 3, 4:
		// OK.
	default:
		panic("bad PLL divider")
	}
	cfgr := rcc.CFGR_Bits(sdiv-1) * rcc.PLLDIV_0
	// Set mul to obtain PLLCLK = 96 MHz (need by USB).
	switch osc {
	case 2:
		cfgr |= rcc.PLLMUL48
	case 3:
		cfgr |= rcc.PLLMUL32
	case 4:
		cfgr |= rcc.PLLMUL24
	case 6:
		cfgr |= rcc.PLLMUL16
	case 8:
		cfgr |= rcc.PLLMUL12
	case 12:
		cfgr |= rcc.PLLMUL8
	case 16, 0:
		cfgr |= rcc.PLLMUL6
	case 24:
		cfgr |= rcc.PLLMUL4
	default:
		panic("bad HSE osc freq")
	}
	// HSE needs milliseconds to stabilize, so enable it now.
	if osc == 0 {
		RCC.HSION().Set()
	} else {
		RCC.HSEON().Set()
	}

	// Setup linear voltage regulator scaling.
	// RCC.PWREN().Set()
	// pwr.PWR.VOS().Store(pwr.VOS_0)
	// RCC.PWREN().Clear()

	// Calculate clock dividers for AHB, APB1, APB2 bus.
	sysclk := 96e6 / uint(sdiv) // Hz
	ahbclk := sysclk
	var apb1clk, apb2clk uint
	if ahbclk <= maxAPBClk {
		cfgr |= rcc.PPRE1_DIV1 | rcc.PPRE2_DIV1
		apb1clk = ahbclk / 1
		apb2clk = ahbclk / 1
	} else {
		cfgr |= rcc.PPRE1_DIV2 | rcc.PPRE2_DIV2
		apb1clk = ahbclk / 2
		apb2clk = ahbclk / 2
	}
	clock[Core] = sysclk
	clock[AHB] = ahbclk
	clock[APB1] = apb1clk
	clock[APB2] = apb2clk

	// Setup Flash.
	FLASH := flash.FLASH
	FLASH.ACC64().Set()
	for FLASH.ACC64().Load() == 0 {
	}
	FLASH.PRFTEN().Set()
	FLASH.LATENCY().Set()

	// Setup PLL.
	if osc == 0 {
		cfgr |= rcc.PLLSRC_HSI
		for RCC.HSIRDY().Load() == 0 {
			// Wait for HSI....
		}
	} else {
		cfgr |= rcc.PLLSRC_HSE
		for RCC.HSERDY().Load() == 0 {
			// Wait for HSE....
		}
	}
	RCC.CFGR.Store(cfgr)
	RCC.PLLON().Set()
	for RCC.PLLRDY().Load() == 0 {
		// Wait for PLL...
	}

	// Change system clock source to PLL.
	RCC.SW().Store(rcc.SW_PLL)
	for RCC.SWS().Load() != rcc.SWS_PLL {
		// Wait for system clock setup...
	}
	RCC.MSION().Clear()
}

// Setup32 wraps SetupPLL to setup MCU to work with 32 MHz clock. See SetupPLL
// for more information.
func Setup32(osc int) {
	SetupPLL(osc, 3)
}
