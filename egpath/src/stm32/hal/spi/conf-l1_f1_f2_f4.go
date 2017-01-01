// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl l1xx_md l1xx_mdp l1xx_hd l1xx_xl f40_41xxx f411xe

package spi

const cr1Mask = ^spi.CR1_Bits(spi.DFF | spi.SPE | spi.BIDIMODE | spi.BIDIOE)

func (p *Periph) setWordSize(size int) {
	if size == 16 {
		p.raw.DFF().Set()
	} else {
		p.raw.DFF().Clear()
	}
}

func (p *Periph) wordSize() int {
	if p.raw.DFF().Load() != 0 {
		return 16
	}
	return 8
}
