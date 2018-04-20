// +build f030x6 f030x

package exti

import (
	"mmio"

	"stm32/hal/raw/syscfg"
)

func exticr(n int) *mmio.U32 {
	return (*mmio.U32)(&syscfg.SYSCFG.EXTICR[n].U32)
}

func exticrEna() {}
func exticrDis() {}
