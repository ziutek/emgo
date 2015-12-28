// +build f10x_hd

package exti

import (
	"stm32/o/f10x_hd/exti"
	"stm32/o/f10x_hd/rcc"
	"stm32/o/f10x_hd/afio"
)

const (
	USB Lines = 1 << 18 // USB wakeup.
)

func extip() *exti.EXTI_Periph { return exti.EXTI }

func exticr(n int) *mmio.U32 {
	return (*mmio.U32)(&afio.AFIO.EXTICR[n].U32)
}
func exticrEna() { rcc.RCC.AFIOEN().Set(); _ = rcc.RCC.APB2ENR.Load() }
func exticrDis() { rcc.RCC.AFIOEN().Clear() }



