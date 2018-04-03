// +build f40_41xxx f411xe f746xx l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package internal

import (
	"bits"
	"unsafe"

	"stm32/hal/raw/rcc"
)

func APB_SetEnabled(addr unsafe.Pointer, en bool) {
	bit := bit(addr, &rcc.RCC.APB1ENR.U32, &rcc.RCC.APB2ENR.U32)
	bit.Store(bits.One(en))
	bit.Load() // Workaround (RCC delay).
}

func APB_Reset(addr unsafe.Pointer) {
	bit := bit(addr, &rcc.RCC.APB1RSTR.U32, &rcc.RCC.APB2RSTR.U32)
	bit.Set()
	bit.Clear()
}

func APB_SetLPEnabled(addr unsafe.Pointer, en bool) {
	bit := bit(addr, &rcc.RCC.APB1LPENR.U32, &rcc.RCC.APB2LPENR.U32)
	bit.Store(bits.One(en))
}
