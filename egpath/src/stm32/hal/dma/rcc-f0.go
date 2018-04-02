// +build  f030x6 f030x8

package dma

import (
	"stm32/hal/raw/rcc"
)

func (p *DMA) enableClock(_ bool) {
	bit := bit(p, &rcc.RCC.AHBENR.U32, rcc.DMAENn)
	bit.Set()
	bit.Load() // RCC delay (workaround for silicon bugs).
}

func (p *DMA) disableClock() {
	bit(p, &rcc.RCC.AHBENR.U32, rcc.DMAENn).Clear()
}

func (p *DMA) reset() {}
