package dma

import (
	"mmio"

	"arch/cortexm/bitband"
)

func bit(p *DMA, reg *mmio.U32) bitband.Bit {
	return bitband.Alias32(reg).Bit(pnum(p))
}
