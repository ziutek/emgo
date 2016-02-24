// +build f40_41xxx f411xe

package dma

import (
	"unsafe"

	"stm32/hal/raw/dma"
	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

type dmaregs struct {
	raw dma.DMA_Periph
	sts [8]dma.DMA_Stream_Periph
}

func enableClock(p *DMA, _ bool) {
	bit := bit(p, &rcc.RCC.AHB1ENR.U32, rcc.DMA1ENn)
	bit.Set()
	bit.Load() // RCC delay (workaround for silicon bugs).
}

func disableClock(p *DMA) {
	bit(p, &rcc.RCC.AHB1ENR.U32, rcc.DMA1ENn).Clear()
}

func reset(p *DMA) {
	bit := bit(p, &rcc.RCC.AHB1RSTR.U32, rcc.DMA1RSTn)
	bit.Set()
	bit.Clear()
}

type chanregs struct {
	raw dma.DMA_Stream_Periph
}

func getChannel(p *DMA, n int) *Channel {
	n--
	if uint(n) > 7 {
		panic(badChan)
	}
	return (*Channel)(unsafe.Pointer(&p.sts[n]))
}
