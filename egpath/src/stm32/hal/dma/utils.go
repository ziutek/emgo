package dma

import (
	"mmio"
	"unsafe"

	"arch/cortexm/bitband"

	"stm32/hal/raw/mmap"
)

const badChan = "dma: bad channel"

// Returns 0 for DMA1, 1 for DMA2.
func dmanum(p *DMA) int {
	return int(uintptr(unsafe.Pointer(p))-mmap.DMA1_BASE) / 0x400
}

func bit(p *DMA, reg *mmio.U32, dma1bitn int) bitband.Bit {
	return bitband.Alias32(reg).Bit(dma1bitn + dmanum(p))
}
