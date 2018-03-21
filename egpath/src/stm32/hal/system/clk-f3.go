// +build f303xe

package system

import "stm32/hal/raw/rcc"

const HSIClk = 8e6 // Hz

const (
	maxAPB1Clk = 36e6 // Hz
)

const (
	ppreDiv1  = rcc.PPRE1_DIV1
	ppreDiv2  = rcc.PPRE1_DIV2
	ppreDiv4  = rcc.PPRE1_DIV4
	ppreDiv8  = rcc.PPRE1_DIV8
	ppreDiv16 = rcc.PPRE1_DIV16

	usbpre = rcc.USBPRE
)
