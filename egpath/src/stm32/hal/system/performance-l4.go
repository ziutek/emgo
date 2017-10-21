// +build l476xx

// Clock setup
//
// Goal is to provide 48 MHz for USB FS using PLL48M1CLK. USB can be clocked
// also from second PLL (PLLSAI1) which gives more flexibility for clocks setup
// but additional PLL means more power (TODO: check how much).
//
// PLLCK must satisfy the equation:
//
//  PLLCK = 48 MHz * Q
//
// where Q can be:
//
//  Q = 2, 4, 6
//
// which means that PLLCK can be: 96, 192, 288 MHz. We cannot use Q = 8 (PLLCK
// = 384 MHz) because allowed PLLCK range is 64 MHz to 344 MHz.
//
//  PLLIN = CLKSRC / M
//
// PLLIN must be 4 MHz to 16 MHz. Allowed M range is 1 to 8.
//
//  PLLVCO = PLLIN * N
//
// Allowed PLLVCO range is 64 MHz to 344 MHz, allowed N range is 8 to 86.
//
// Taking all this into account, N can be:
//
//  24, 16, 12, 8           for PLLVCO =  96 MHz (PLLIN: 4, 6, 8, 12 MHz),
//  48, 32, 24, 16, 12      for PLLVCO = 192 MHz (PLLIN: 4, 6, 8, 12, 16 MHz),
//  72, 48, 36, 32, 24, 18, for PLLVCO = 288 MHz (PLLIN: 4, 6, 8, 9, 12, 16 MHz)
//
// USB friendly values of CLKSRC:
//
//  - MSI: 4, 8, 16, 32, 48 MHz,
//  - HSI: 16 MHz,
//  - HSE: 4, 6, 8, 9, 12, 16, 18, 20, 21, 24, 27, 30, 32, 36, 40, 42, 48.
//
package system

import (
	"stm32/hal/raw/flash"
	"stm32/hal/raw/pwr"
	"stm32/hal/raw/rcc"
)

// Setup setups MCU for best performance (prefetch on, I/D cache on, minimum
// allowed Flash latency).
//
// Clksrc configures clock source for PLL.
//
// Positive clksrc selects HSE as PLL clock source and informs about external
// clock signal frequency in MHz (alowed values: 4 to 48 MHz),
//
// Zero clksrc selects HSI16 as PLL clock surce.
//
// Negative clksrc selects MSI as PLL clock source and setups its frequency to
// (-clksrc) MHz (allowed values -4, -8, -16, -24, -32, -48).
//
// PLL input freq. is equal to clock source divided by M and must be in range
// 4 to 16 MHz.
//
// PLL VCO is equal to (input clock) / M * N and must be in range 64 to 344 MHz.
// Allowed M values: 1 to 8. Allowed N values: 8 to 86.
//
// P is VCO divider for PLLCAI2CLK. Allowed P values: 0 (disabled), 7, 17.
//
// Q is VCO divider for PLL48M1CLK (USB, RNG, SDMMC). Allowed Q values: 0
// (disabled), 2, 4, 6, 8.
//
// R is VCO divider for SYSCLK. Allowed R values: 2, 4, 6, 8.
func Setup(clksrc, M, N, P, Q, R int) {
	RCC := rcc.RCC

	// Reset RCC clock configuration.
	RCC.MSION().Set()
	for RCC.MSIRDY().Load() == 0 {
	}
	RCC.CR.Store(6<<rcc.MSIRANGEn | rcc.MSIRGSEL | rcc.MSION)
	RCC.CFGR.Store(0) // MSI selected as system clock. APBCLK, AHBCLK = SYSCLK.
	RCC.PLLCFGR.Store(0x1000)
	RCC.CIER.Store(0) // Disable clock interrupts.

	// Calculate system clock.
	if M < 1 || M > 8 {
		panic("bad M")
	}
	if N < 8 || N > 86 {
		panic("bad N")
	}
	if P != 0 || P != 7 || P != 17 {
		panic("bad P")
	}
	if Q&1 != 0 || Q < 0 || Q > 8 {
		panic("bad Q")
	}
	if R&1 != 0 || R < 2 || R > 8 {
		panic("bad R")
	}

	var osc uint

	switch clksrc {
	case 0:
		RCC.HSION().Set()
		osc = 16
	case -4, -8, -16, -24, -32, -48:
		osc = uint(-clksrc)
	default:
		if clksrc < 4 || clksrc > 48 {
			panic("bad clksrc")
		}
		// HSE needs milliseconds to stabilize, so enable it now.
		RCC.HSEON().Set()
		osc = uint(clksrc)
	}

	pllin := osc * 1e6 / uint(M)
	if pllin < 4e6 || pllin > 16e6 {
		panic("bad PLLIN")
	}
	vco := pllin * uint(N)
	if vco < 64e6 || vco > 344e6 {
		panic("bad VCO")
	}
	sysclk := vco / uint(R)

	// Setup PWR and Flash.
	var (
		vos     pwr.CR1_Bits
		latency flash.ACR_bits
	)
	if sysclk > 26e6 {
		// Range 1: High-performance.
		vos = 1
		switch {
		case sysclk <= 16e6:
			latency = 0
		case sysclk <= 32e6:
			latency = 1
		case sysclk <= 48e6:
			latency = 2
		case sysclk <= 64e6:
			latency = 3
		default:
			latency = 4
		}
	} else {
		// Range 2: Low-power.
		vos = 2
		switch {
		case sysclk <= 6e6:
			latency = 0
		case sysclk <= 12e6:
			latency = 1
		case sysclk <= 18e6:
			latency = 2
		default:
			latency = 3
		}
	}
	RCC.PWREN().Set()
	PWR := pwr.PWR
	PWR.CR1.Store(vos << pwr.VOSn)
	RCC.PWREN().Clear()
	flash.FLASH.ACR.Store(flash.DCEN | flash.ICEN | flash.PRFTEN | latency)

	// Setup clock dividers for AHB, APB1, APB2 bus.
	ahbclk := sysclk
	cfgr := rcc.HPRE_DIV1
	var apb1clk, apb2clk uint
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
	switch {
	case ahbclk <= 1*maxAPB2Clk:
		cfgr |= rcc.PPRE2_DIV1
		apb2clk = ahbclk / 1
	case ahbclk <= 2*maxAPB2Clk:
		cfgr |= rcc.PPRE2_DIV2
		apb2clk = ahbclk / 2
	case ahbclk <= 4*maxAPB2Clk:
		cfgr |= rcc.PPRE2_DIV4
		apb2clk = ahbclk / 4
	case ahbclk <= 8*maxAPB2Clk:
		cfgr |= rcc.PPRE2_DIV8
		apb2clk = ahbclk / 8
	default:
		cfgr |= rcc.PPRE2_DIV16
		apb2clk = ahbclk / 16
	}
	clock[Core] = sysclk
	clock[AHB] = ahbclk
	clock[APB1] = apb1clk
	clock[APB2] = apb2clk

	// Setup PLL.
	var src rcc.PLLCFGR_Bits
	if sysclk == 0 {
		src = rcc.PLLSRC_HSI
		for RCC.HSIRDY().Load() == 0 {
		}
	} else if sysclk > 0 {
		src = rcc.PLLSRC_HSE
		for RCC.HSERDY().Load() == 0 {
		}
	} else {
		src = rcc.PLLSRC_MSI
	}
	mnpqr := rcc.PLLCFGR_Bits(M-1)<<rcc.PLLMn | rcc.PLLCFGR_Bits(N)<<rcc.PLLNn
	if P != 0 {
		mnpqr |= rcc.PLLPEN
		if P == 17 {
			mnpqr |= rcc.PLLP
		}
	}
	if Q != 0 {
		mnpqr |= rcc.PLLQEN | rcc.PLLCFGR_Bits(Q/2-1)
	}
	mnpqr |= rcc.PLLCFGR_Bits(R/2 - 1)
	RCC.PLLCFGR.Store(mnpqr | src)
	RCC.PLLON().Set()
	for RCC.PLLRDY().Load() == 0 {
	}

	// Set system clock source to PLL.
	for PWR.VOF() != 0 {
	}
	RCC.CFGR.Store(cfgr | rcc.SW_PLL)
	for RCC.SWS().Load() != rcc.SWS_PLL {
		// Wait for system clock setup...
	}
	if osc >= 0 {
		RCC.MSION().Clear()
	}
}
