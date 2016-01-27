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
	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

// Setup setups MCU for best performance (Flash prefetch and 64-bit access
// on).
//
// osc is freqency of external resonator in MHz. Allowed values: 2, 3, 4, 6, 8,
// 12, 16, 24. Use 0 to select internal HSI oscilator as system clock source.
//
// sdiv is system clock divider. Allowed values: 2, 3, 4. sdiv determine the
// system clock frequency according to the formula:
//
//  SysClk = 96e6 / sdiv [Hz]
//
func Setup(osc, sdiv int) {
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

	SysClk = 96e6 / uint(sdiv) // Hz

	// Setup linear voltage regulator scaling.
	// RCC.PWREN().Set()
	// pwr.PWR.VOS().Store(pwr.VOS_0)
	// RCC.PWREN().Clear()

	// Calculate clock dividers for AHB, APB1, APB2 bus.
	AHBClk = SysClk
	if AHBClk <= 32e6 {
		cfgr |= rcc.PPRE1_DIV1 | rcc.PPRE2_DIV1
		APB1Clk = AHBClk / 1
		APB2Clk = AHBClk / 1
	} else {
		cfgr |= rcc.PPRE1_DIV2 | rcc.PPRE2_DIV2
		APB1Clk = AHBClk / 2
		APB2Clk = AHBClk / 2
	}

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

	setupOS()
}

// Setup32 setups MCU to work with 96 MHz clock.
// See Setup for description of osc.
func Setup32(osc int) {
	Setup(osc, 3)
}

func PeriphClk(baseaddr uintptr) uint {
	switch {
	case baseaddr >= mmap.AHBPERIPH_BASE:
		return AHBClk
	case baseaddr >= mmap.APB2PERIPH_BASE:
		return APB2Clk
	case baseaddr >= mmap.APB1PERIPH_BASE:
		return APB1Clk
	}
	return 0
}
