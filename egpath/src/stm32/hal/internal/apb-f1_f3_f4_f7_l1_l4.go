// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl f303xe f40_41xxx f411xe f746xx l1xx_md l1xx_mdp l1xx_hd l1xx_xl l476xx

package internal

import (
	"mmio"
	"unsafe"

	"stm32/hal/raw/mmap"
)

func bit(addr unsafe.Pointer, apb1reg, apb2reg *mmio.U32) AtomicBit {
	a := uintptr(addr)
	var reg *mmio.U32
	if a >= mmap.APB2PERIPH_BASE {
		reg = apb2reg
		a -= mmap.APB2PERIPH_BASE
	} else {
		reg = apb1reg
		a -= mmap.APB1PERIPH_BASE
	}
	n := int(a / 0x400)
	return AtomicBit{reg, n}
}
