// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package dma

import (
	"unsafe"

	"stm32/hal/raw/dma"
	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

var (
	DMA1 = DMA{(*registers)(unsafe.Pointer(mmap.DMA1_BASE))}
	DMA2 = DMA{(*registers)(unsafe.Pointer(mmap.DMA2_BASE))}
)

type registers struct {
	dma.DMA_Periph
	chs [7]struct {
		dma.DMA_Channel_Periph
		_ uint32
	}
}

func pnum(p DMA) int {
	return int(uintptr(unsafe.Pointer(p.registers))-mmap.AHBPERIPH_BASE) / 0x400
}

func enableClock(p DMA, _ bool) {
	bit := bit(p, &rcc.RCC.AHBENR.U32)
	bit.Set()
	bit.Load() // RCC delay (workaround for silicon bugs).
}

func disableClock(p DMA) {
	bit(p, &rcc.RCC.AHBENR.U32).Clear()
}


type channel struct {
	p *dma.DMA_Periph
	ch *dma.DMA_Channel_Periph
}