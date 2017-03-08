// +build l476xx

package internal

import (
	"bits"
	"unsafe"

	"stm32/hal/raw/rcc"
)

// BUG: APB1SMENR2 (LPUART1, SWPMI1, LPTIM2) not supported.

func APB_SetEnabled(addr unsafe.Pointer, en bool) {
	bit := bit(addr, &rcc.RCC.APB1ENR1.U32, &rcc.RCC.APB2ENR.U32)
	bit.Store(bits.One(en))
	bit.Load() // Workaround (RCC delay).
}

func APB_Reset(addr unsafe.Pointer) {
	bit := bit(addr, &rcc.RCC.APB1RSTR1.U32, &rcc.RCC.APB2RSTR.U32)
	bit.Set()
	bit.Clear()
}

func APB_SetLPEnabled(addr unsafe.Pointer, en bool) {
	bit := bit(addr, &rcc.RCC.APB1SMENR1.U32, &rcc.RCC.APB2SMENR.U32)
	bit.Store(bits.One(en))
}
