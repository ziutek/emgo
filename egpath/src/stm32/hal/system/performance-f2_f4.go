// +build f40_41xxx f411xe f2xxx-TODO

// Clock setup
//
// Goal is to provide 48 MHz for USB. So PLLCK must satisfy the equation:
//
//  PLLCK = 48 MHz * Q
//
// where Q = 2..15, which means that PLLCK can be: 96, 144, ... , 720 MHz.
//
// But allowed PLLCK is between 100 and 432 MHz so useful Q values are:
//
//  Q = 3..9
//
// which means PLLCK can be: 144, 192, 240, 288, 336, 384, 432 MHz.
//
// PLL multipler N range is 50..432. There is recommendation to use 2 MHz input
// clock to PLL to limit its jitter. Taking this into account PCLK can be:
//
//  PLLCK = N * 2 MHz
//
// PCLK should be between 100 and 432 MHz so useful N values are:
//
//  N = 50..216
//
// There is much smaller choice of N values that satisfy USB requirements:
//
//  N = 72, 96, 120, 144, 168, 192, 216.
//
// System Clock is derived from PLLCK as follows:
//
//  SYSCLK = PLLCK / P
//
// where P = 2, 4, 6, 8.
//
// If 2 MHz PLL input clock is used, System Clock can be calculated as follows:
//
//  SYSCLK = 2 MHz * N / P
//
package system

import (
	"stm32/hal/raw/flash"
	"stm32/hal/raw/rcc"
)

// Setup setups MCU for best performance (prefetch on, I/D cache on, minimum
// allowed Flash latency).
//
// osc is freqency of external resonator in MHz. Allowed values are multiples
// of 2, from 4 to 26. Use 0 to select internal HSI oscilator as system clock
// source.
//
// mul is PLL multipler. Allowed values are from 50 to 216 but if USB will be
// used, mul can be only:
//
//  mul(USB) = 72, 96, 120, 144, 168, 192, 216
//
// sdiv is system clock divider. Allowed values: 2, 4, 6, 8.
//
// Both mul and sdiv determine the system clock frequency according to the
// formula:
//
//  SysClk = 2e6 * mul / sdiv [Hz]
//
func Setup(osc, mul, sdiv int) {
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
		panic("bad HSE osc freq")
	}
	if mul < 72 || mul > 216 {
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

	// Setup linear voltage regulator scaling.
	// RCC.PWREN().Set()
	// pwr.PWR.VOS().Store(pwr.VOS_1 | pwr.VOS_0)
	// RCC.PWREN().Clear()

	// Setup clock dividers for AHB, APB1, APB2 bus.
	sysclk := 2e6 * uint(mul) / uint(sdiv)
	ahbclk := sysclk
	var apb1clk, apb2clk uint
	switch {
	case ahbclk <= 1*maxAPB1Clk:
		apb1clk = ahbclk / 1
	case ahbclk <= 2*maxAPB1Clk:
		RCC.PPRE1().Store(rcc.PPRE1_DIV2)
		apb1clk = ahbclk / 2
	case ahbclk <= 4*maxAPB1Clk:
		RCC.PPRE1().Store(rcc.PPRE1_DIV4)
		apb1clk = ahbclk / 4
	case ahbclk <= 8*maxAPB1Clk:
		RCC.PPRE1().Store(rcc.PPRE1_DIV8)
		apb1clk = ahbclk / 8
	default:
		RCC.PPRE1().Store(rcc.PPRE1_DIV16)
		apb1clk = ahbclk / 16
	}
	switch {
	case ahbclk <= 1*maxAPB2Clk:
		apb2clk = ahbclk / 1
	case ahbclk <= 2*maxAPB2Clk:
		RCC.PPRE2().Store(rcc.PPRE2_DIV2)
		apb2clk = ahbclk / 2
	case ahbclk <= 4*maxAPB2Clk:
		RCC.PPRE2().Store(rcc.PPRE2_DIV4)
		apb2clk = ahbclk / 4
	case ahbclk <= 8*maxAPB2Clk:
		RCC.PPRE2().Store(rcc.PPRE2_DIV8)
		apb2clk = ahbclk / 8
	default:
		RCC.PPRE2().Store(rcc.PPRE2_DIV16)
		apb2clk = ahbclk / 16
	}
	clock[Core] = sysclk
	clock[AHB] = ahbclk
	clock[APB1] = apb1clk
	clock[APB2] = apb2clk

	// Setup Flash.
	FLASH := flash.FLASH
	latency := flash.ACR_Bits((sysclk-1)/30e6) * flash.LATENCY_1WS
	FLASH.ACR.SetBits(flash.DCEN | flash.ICEN | flash.PRFTEN | latency)

	// Setup PLL.
	var (
		src rcc.PLLCFGR_Bits
		M   rcc.PLLCFGR_Bits                               // PLL input divider.
		N   = rcc.PLLCFGR_Bits(mul) * rcc.PLLN_0           // PLL multiler.
		P   = rcc.PLLCFGR_Bits(sdiv/2-1) * rcc.PLLP_0      // SysClk divider.
		Q   = rcc.PLLCFGR_Bits(2*mul+47) / 48 * rcc.PLLQ_0 // USB 48MHz divider.
	)
	if osc != 0 {
		src = rcc.PLLSRC_HSE
		M = rcc.PLLCFGR_Bits(osc/2) * rcc.PLLM_0
		for RCC.HSERDY().Load() == 0 {
			// Wait for HSE...
		}
	} else {
		src = rcc.PLLSRC_HSI
		M = 16 / 2 * rcc.PLLM_0
	}
	RCC.PLLCFGR.Store(Q | src | P | N | M)
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

// Setup96 setups MCU to work with 96 MHz clock.
// See Performance for description of osc.
func Setup96(osc int) {
	Setup(osc, 192, 4)
}

// Setup168 setups MCU to work with 168 MHz clock.
// See Performance for description of osc.
func Setup168(osc int) {
	Setup(osc, 168, 2)
}