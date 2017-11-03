// +build f40_41xxx f411xe f746xx

package dma

import (
	"bits"

	"stm32/hal/raw/rcc"
)

func (p *DMA) enableClock(lp bool) {
	enbit := bit(p, &rcc.RCC.AHB1ENR.U32, rcc.DMA1ENn)
	enbit.Set()
	bit(p, &rcc.RCC.AHB1LPENR.U32, rcc.DMA1LPENn).Store(bits.One(lp))
	enbit.Load() // RCC delay (workaround for silicon bugs).
}

func (p *DMA) disableClock() {
	bit(p, &rcc.RCC.AHB1ENR.U32, rcc.DMA1ENn).Clear()
}

func (p *DMA) reset() {
	bit := bit(p, &rcc.RCC.AHB1RSTR.U32, rcc.DMA1RSTn)
	bit.Set()
	bit.Clear()
}
