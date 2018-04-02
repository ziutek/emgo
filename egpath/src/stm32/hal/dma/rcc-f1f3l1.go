// +build  f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl f303xe l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package dma

import (
	"stm32/hal/raw/rcc"
)

func (p *DMA) enableClock(_ bool) {
	bit := bit(p, &rcc.RCC.AHBENR.U32, rcc.DMA1ENn)
	bit.Set()
	bit.Load() // RCC delay (workaround for silicon bugs).
}

func (p *DMA) disableClock() {
	bit(p, &rcc.RCC.AHBENR.U32, rcc.DMA1ENn).Clear()
}

func (p *DMA) reset() {}
