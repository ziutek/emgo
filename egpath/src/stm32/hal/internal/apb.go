package internal

import (
	"bits"
	"mmio"
	"unsafe"

	"arch/cortexm/bitband"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
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

func APB_SetEnabled(addr unsafe.Pointer, en bool) {
	bit := bit(addr, &rcc.RCC.APB1ENR.U32, &rcc.RCC.APB2ENR.U32)
	bit.Store(bits.One(en))
	bit.Load() // Workaround (RCC delay).
}

func APB_Reset(addr unsafe.Pointer) {
	bit := bit(addr, &rcc.RCC.APB1RSTR.U32, &rcc.RCC.APB2RSTR.U32)
	bit.Set()
	bit.Clear()
}
