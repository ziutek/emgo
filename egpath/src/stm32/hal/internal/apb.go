//  +build !cortexm7f !cortexm7d
package internal

import (
	"mmio"
	"unsafe"

	"arch/cortexm/bitband"

	"stm32/hal/raw/mmap"
)

func bit(addr unsafe.Pointer, apb1reg, apb2reg *mmio.U32) bitband.Bit {
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
	return bitband.Alias32(reg).Bit(n)
}
