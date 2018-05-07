// +build f030x6 f030x8

package system

import "stm32/hal/raw/rcc"

const HSIClk = 8e6 // Hz

const (
	maxAPB1Clk = 48e6 // Hz
)

const (
	ppreDiv1  = rcc.PPRE_DIV1
	ppreDiv2  = rcc.PPRE_DIV2
	ppreDiv4  = rcc.PPRE_DIV4
	ppreDiv8  = rcc.PPRE_DIV8
	ppreDiv16 = rcc.PPRE_DIV16

	usbpre = 0
)
