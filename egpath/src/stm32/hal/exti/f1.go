// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package exti

import (
	"mmio"

	"stm32/hal/raw/afio"
	"stm32/hal/raw/exti"
	"stm32/hal/raw/rcc"
)

const (
	USB   Lines = 1 << 18 // USB wakeup.
	Ether Lines = 1 << 19 // Ethernet wakeup.
)

func extip() *exti.EXTI_Periph { return exti.EXTI }

func exticr(n int) *mmio.U32 {
	return (*mmio.U32)(&afio.AFIO.EXTICR[n].U32)
}
func exticrEna() { rcc.RCC.AFIOEN().Set(); _ = rcc.RCC.APB2ENR.Load() }
func exticrDis() { rcc.RCC.AFIOEN().Clear() }
