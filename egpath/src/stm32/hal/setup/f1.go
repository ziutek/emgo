// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl

package setup

import (
	"stm32/hal/raw/flash"
	"stm32/hal/raw/rcc"
)

// Performance setups MCU for best performance ().
//
// osc is freqency of external resonator in MHz. Allowed values: 4 to 16 MHz.
// Use 0 to select internal HSI oscilator (8 MHz / 2) as system clock source.
//
// sdiv is system clock divider.
//
// div2 and mul determine the system clock frequency according to the formula:
//
//  SysClk = osc / (div2 ? 2 : 1) * mul * 1e6 [Hz]
//
// when osc != 0 or:
//
//  SysClk = 4e6 * mul [Hz]
//
// when osc == 0. mul can be 2..16.
//
// USB requires HSE and SysClk set to 48e6 or 72e6 Hz.
func Performance(osc, mul int, div2 bool) {
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
	if osc != 0 && (osc < 4 || osc > 16) {
		panic("bad HSE osc freq")
	}
	if mul < 2 || mul > 16 {
		panic("bad PLL multipler")
	}
	if osc == 0 {
		SysClk = 4e6 * uint(mul) // Hz
	} else {
		// HSE needs milliseconds to stabilize, so enable it now.
		RCC.HSEON().Set()
		SysClk = uint(osc) * uint(mul) * 1e6 // Hz
		if div2 {
			SysClk /= 2
		}
	}
	AHBClk = SysClk
	APB2Clk = AHBClk
	cfgr := rcc.CFGR_Bits(mul-2) * rcc.PLLMULL_0
	const maxAPB1Clk = 36e6 // Hz
	switch {
	case AHBClk <= 1*maxAPB1Clk:
		cfgr |= rcc.PPRE1_DIV1
		APB1Clk = AHBClk / 1
	case AHBClk <= 2*maxAPB1Clk:
		cfgr |= rcc.PPRE1_DIV2
		APB1Clk = AHBClk / 2
	case AHBClk <= 4*maxAPB1Clk:
		cfgr |= rcc.PPRE1_DIV4
		APB1Clk = AHBClk / 4
	case AHBClk <= 8*maxAPB1Clk:
		cfgr |= rcc.PPRE1_DIV8
		APB1Clk = AHBClk / 8
	default:
		cfgr |= rcc.PPRE1_DIV16
		APB1Clk = AHBClk / 16
	}
	if SysClk <= 48e6 {
		cfgr |= rcc.USBPRE
	}
	// Setup Flash.
	FLASH := flash.FLASH
	latency := flash.ACR_Bits((SysClk-1)/24e6) * flash.LATENCY_1
	FLASH.ACR.SetBits(flash.PRFTBE | latency)

	// Setup PLL.
	if osc == 0 {
		cfgr |= rcc.PLLSRC_HSI_Div2
	} else {
		cfgr |= rcc.PLLSRC_HSE
		if div2 {
			cfgr |= rcc.PLLXTPRE_HSE_Div2
		}
		for RCC.HSERDY().Load() == 0 {
			// Wait for HSE...
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
	if osc != 0 {
		RCC.HSION().Clear()
	}

	setupOS()
}
