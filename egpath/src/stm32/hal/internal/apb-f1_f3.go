// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl f303xe

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

func APB_SetLPEnabled(_ unsafe.Pointer, _ bool) {}
