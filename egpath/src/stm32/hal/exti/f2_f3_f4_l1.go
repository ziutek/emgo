// +build f40_41xxx f411xe l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package exti

import (
	"mmio"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/syscfg"
)

func exticr(n int) *mmio.U32 {
	return (*mmio.U32)(&syscfg.SYSCFG.EXTICR[n].U32)
}
func exticrEna() { rcc.RCC.SYSCFGEN().Set(); _ = rcc.RCC.APB2ENR.Load() }
func exticrDis() { rcc.RCC.SYSCFGEN().Clear() }
