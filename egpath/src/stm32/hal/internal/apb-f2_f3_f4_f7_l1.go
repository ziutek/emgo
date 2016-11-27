// +build f40_41xxx f411xe f746xx l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package internal

import (
	"bits"
	"unsafe"

	"stm32/hal/raw/rcc"
)

func APB_SetLPEnabled(addr unsafe.Pointer, en bool) {
	bit := bit(addr, &rcc.RCC.APB1LPENR.U32, &rcc.RCC.APB2LPENR.U32)
	bit.Store(bits.One(en))
}
