// +build l476xx

package dma

import (
	"stm32/hal/raw/rcc"
)

func (p *DMA) enableClock(_ bool) {
	bit := bit(p, &rcc.RCC.AHB1ENR.U32, rcc.DMA1ENn)
	bit.Set()
	bit.Load() // RCC delay (workaround for silicon bugs).
}

func (p *DMA) disableClock() {
	bit(p, &rcc.RCC.AHB1ENR.U32, rcc.DMA1ENn).Clear()
}

func (p *DMA) reset() {
	bit := bit(p, &rcc.RCC.AHB1RSTR.U32, rcc.DMA1RSTn)
	bit.Set()
	bit.Clear()
}
