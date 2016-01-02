// +build l1xx_md l1xx_mdp l1xx_hd l1xx_xl

// Package setup allows to easy setup MCU for typical use.
//
// Clock setup
//
// Goal is to provide 48 MHz for USB so PLLCLK must be set to 96 MHz because the
// USBCLK = PLLCLK / 2.
//
// System Clock is derived from PLLCK as follows:
//
//  SYSCLK = PLLCK / PLLDIV
//
package setup

import (
	"stm32/hal/raw/flash"
	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

var (
	SysClk  uint // System clock [Hz]
	AHBClk  uint // AHB clock [Hz]
	APB1Clk uint // APB1 clock [Hz]
	APB2Clk uint // APB2 clock [Hz]
)

// Performance setups MCU for best performance (Flash prefetch and 64-bit access
// on).
//
// osc is freqency of external resonator in MHz. Allowed values: 2, 3, 4, 6, 8,
// 12, 16, 24. Use 0 to select internal HSI oscilator as system clock source.
//
// sdiv is system clock divider. Allowed values: 2, 3, 4. sdiv determine the
// system clock frequency according to the formula:
//
//  SysClk = 96 MHz / sdiv
//
func Performance(osc, sdiv int) {
	RCC := rcc.RCC

	// Reset RCC clock configuration.
	RCC.MSION().Set()
	for RCC.MSIRDY().Load() == 0 {
		// Wait for MSI...
	}
	RCC.CFGR.Store(0)
	for RCC.SWS().Load() != rcc.SWS_MSI {
		// Wait for system clock setup...
	}
	RCC.CR.ClearBits(rcc.HSION | rcc.HSEON | rcc.CSSON | rcc.PLLON | rcc.HSEBYP)
	RCC.CIR.Store(0) // Disable clock interrupts.

	// Set mul to obtain PLLCLK = 96 MHz (need by USB).
	var mul rcc.CFGR_Bits
	switch osc {
	case 2:
		mul = rcc.PLLMUL48
	case 3:
		mul = rcc.PLLMUL32
	case 4:
		mul = rcc.PLLMUL24
	case 6:
		mul = rcc.PLLMUL16
	case 8:
		mul = rcc.PLLMUL12
	case 12:
		mul = rcc.PLLMUL8
	case 16, 0:
		mul = rcc.PLLMUL6
	case 24:
		mul = rcc.PLLMUL4
	default:
		panic("bad HSE osc freq")
	}
	switch sdiv {
	case 2, 3, 4:
		// OK.
	default:
		panic("bad PLL divider")
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

	// Disable AHB clock (if enabled before).
	RCC.HPRE().Store(0)

	AHBClk = SysClk
	ahbdiv := rcc.HPRE_DIV1
	if AHBClk <= 32e6 {
		RCC.PPRE1().Store(rcc.PPRE1_DIV1)
		RCC.PPRE2().Store(rcc.PPRE2_DIV1)
		APB1Clk = AHBClk / 1
		APB2Clk = AHBClk / 1
	} else {
		RCC.PPRE1().Store(rcc.PPRE1_DIV2)
		RCC.PPRE2().Store(rcc.PPRE2_DIV2)
		APB1Clk = AHBClk / 2
		APB2Clk = AHBClk / 2
	}

	// Enable AHB clock.
	RCC.HPRE().Store(ahbdiv)

	// Setup Flash.
	FLASH := flash.FLASH
	FLASH.ACC64().Set()
	for FLASH.ACC64().Load() == 0 {
	}
	FLASH.PRFTEN().Set()
	FLASH.LATENCY().Set()

	// Setup PLL.
	if osc == 0 {
		for RCC.HSIRDY().Load() == 0 {
			// Wait for HSI....
		}
		RCC.PLLSRC().Store(rcc.PLLSRC_HSI)
	} else {
		for RCC.HSERDY().Load() == 0 {
			// Wait for HSE....
		}
		RCC.PLLSRC().Store(rcc.PLLSRC_HSE)
	}
	RCC.CFGR.StoreBits(
		rcc.PLLDIV|rcc.PLLMUL,
		rcc.CFGR_Bits(sdiv-1)*rcc.PLLDIV_0|mul,
	)
	RCC.PLLON().Set()

	for FLASH.LATENCY().Load() != flash.LATENCY {
		// Ensure flash latency is set before incrase frequency.
	}
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

// Performance32 setups MCU to work with 96 MHz clock.
// See Performance for description of osc.
func Performance32(osc int) {
	Performance(osc, 3)
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
