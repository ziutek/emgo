// +build f030x6 f030x8

package internal

import (
	"mmio"
	"unsafe"

	"stm32/hal/raw/mmap"
)

func bit(addr unsafe.Pointer, apb1reg, apb2reg *mmio.U32) AtomicBit {
	a := uintptr(addr)
	var reg *mmio.U32
	if a >= mmap.APBPERIPH_BASE+0x10000 {
		reg = apb2reg
		a -= mmap.APBPERIPH_BASE + 0x10000
	} else {
		reg = apb1reg
		a -= mmap.APBPERIPH_BASE
	}
	n := int(a / 0x400)
	return AtomicBit{reg, n}
}
