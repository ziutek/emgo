// +build f303xe l476xx

package spi

import (
	"stm32/hal/raw/spi"
)

const cr1Mask = ^spi.CR1_Bits(spi.SPE | spi.BIDIMODE | spi.BIDIOE)

func (p *Periph) setWordSize(size int) {
	ds := spi.CR2_Bits((size - 1) & 0xf << spi.DSn)
	if size <= 8 {
		ds |= spi.FRXTH
	}
	p.raw.CR2.StoreBits(spi.FRXTH|spi.DS, ds)
}

func (p *Periph) wordSize() int {
	return int(p.raw.CR2.Bits(spi.DS))>>spi.DSn&0xf + 1
}
