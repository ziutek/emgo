package setup

import (
	"stm32/o/f40_41xxx/flash"
	"stm32/o/f40_41xxx/rcc"
)

var (
	SysClk  uint // System clock [Hz]
	AHBClk  uint // AHB clock [Hz]
	APB1Clk uint // APB1 clock [Hz]
	APB2Clk uint // APB2 clock [Hz]
)

// Performance setups MCU for best performance (prefetch on, I/D cache on,
// minimum allowed Flash latency).
//
// osc is freqency of external resonator in MHz. Allowed values are multiples
// of 2, from 4 to 26. Use 0 to select internal HSI oscilator as system clock
// source.
//
// mul is PLL multipler. Allowed values are multiples of 24, from 96 to 216.
//
// sdiv is system clock divider. Allowed values: 2, 4, 6, 8.
//
// Both mul and sdiv determine the system clock frequency according to the
// formula:
//
//  SysClk = 2 MHz * mul / sdiv
//
// Performance ensures that pheripheral clock used by SDIO, RNG and USBFS is
// equal to 48 MHz (required by USB FS).
func Performance(osc, mul, sdiv int) {
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
	RCC.PLLCFGR.Store(0x24003010)
	RCC.CIR.Store(0) // Disable clock interrupts.

	// Calculate system clock.
	if osc != 0 && (osc < 4 || osc > 26 || osc&1 != 0) {
		panic("bad HSE freq")
	}
	if mul < 96 || mul > 216 || mul%24 != 0 {
		panic("bad PLL N multipler")
	}
	switch sdiv {
	case 2, 4, 6, 8:
		// OK.
	default:
		panic("bad PLL P divider")
	}
	// HSE needs milliseconds to stabilize, so enable it now.
	if osc != 0 {
		RCC.HSEON().Set()
	}
	sysclk := 2 * mul / sdiv    // MHz
	SysClk = uint(sysclk) * 1e6 // Hz

	// Setup linear voltage regulator scaling.
	// RCC.PWREN().Set()
	// pwr.PWR.VOS().Store(pwr.VOS_1 | pwr.VOS_0)
	// RCC.PWREN().Clear()

	// Calculate clock dividers for AHB, APB1, APB2 bus.
	// APB1Freq must be <= 42 MHz.
	// APB2Freq must be <= 84 MHz.
	AHBClk = SysClk
	var apb1div, apb2div rcc.CFGR_Bits
	switch {
	case sysclk <= 42:
		apb1div = rcc.PPRE1_DIV1
		apb2div = rcc.PPRE2_DIV1
		APB1Clk = SysClk
		APB2Clk = SysClk
	case sysclk <= 84:
		apb1div = rcc.PPRE1_DIV2
		apb2div = rcc.PPRE2_DIV1
		APB1Clk = SysClk / 2
		APB2Clk = SysClk
	case sysclk <= 168:
		apb1div = rcc.PPRE1_DIV4
		apb2div = rcc.PPRE2_DIV2
		APB1Clk = SysClk / 4
		APB2Clk = SysClk / 2
	default:
		apb1div = rcc.PPRE1_DIV8
		apb2div = rcc.PPRE2_DIV4
		APB1Clk = SysClk / 8
		APB2Clk = SysClk / 4
	}
	// Set clock dividers for AHB, APB1, APB2 bus.
	RCC.HPRE().Store(rcc.HPRE_DIV1)
	RCC.PPRE1().Store(apb1div)
	RCC.PPRE2().Store(apb2div)

	// Setup Flash.
	FLASH := flash.FLASH
	latency := flash.ACR_Bits((sysclk-1)/30) * flash.LATENCY_1WS
	FLASH.ACR.SetBits(flash.DCEN | flash.ICEN | flash.PRFTEN | latency)

	// Setup PLL.
	var (
		pllsrc rcc.PLLCFGR_Bits
		M      rcc.PLLCFGR_Bits                          // PLL input divider.
		N      = rcc.PLLCFGR_Bits(mul) * rcc.PLLN_0      // PLL multiler.
		P      = rcc.PLLCFGR_Bits(sdiv/2-1) * rcc.PLLP_0 // SysClk divider
		Q      = rcc.PLLCFGR_Bits(mul/24) * rcc.PLLQ_0   // USB 48 MHz divider.
	)
	if osc != 0 {
		pllsrc = rcc.PLLSRC_HSE
		M = rcc.PLLCFGR_Bits(osc/2) * rcc.PLLM_0
		for RCC.HSERDY().Load() == 0 {
			// Wait for HSE....
		}
	} else {
		pllsrc = rcc.PLLSRC_HSI
		M = 16 / 2 * rcc.PLLM_0
	}
	RCC.PLLCFGR.Store(Q | pllsrc | P | N | M)
	RCC.PLLON().Set()

	for FLASH.LATENCY().Load() != latency {
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
	if osc != 0 {
		RCC.HSION().Clear()
	}

	setupOS()
}

// Performance84 setups MCU to work with 84 MHz clock.
// See Performance for description of osc.
func Performance84(osc int) {
	Performance(osc, 168, 4)
}

// Performance168 setups MCU to work with 168 MHz clock.
// See Performance for description of osc.
func Performance168(osc int) {
	Performance(osc, 168, 2)
}
