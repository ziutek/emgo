package dma

import (
	"mmio"
	"unsafe"

	"stm32/hal/internal"
	"stm32/hal/raw/mmap"
)

const badStream = "dma: bad stream"

// Returns 0 for DMA1, 1 for DMA2.
func dmanum(p *DMA) int {
	return int(uintptr(unsafe.Pointer(p))-mmap.DMA1_BASE) / 0x400
}

func bit(p *DMA, r *mmio.U32, dma1bitn int) internal.AtomicBit {
	return internal.AtomicBit{r, dma1bitn + dmanum(p)}
}
