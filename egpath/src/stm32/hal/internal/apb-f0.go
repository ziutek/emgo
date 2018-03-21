// +build f030x6

package internal

import (
	"bits"
	"mmio"
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

func bit(addr unsafe.Pointer, apb1reg, _ *mmio.U32) AtomicBit {
	n := int(uintptr(addr)-mmap.APBPERIPH_BASE) / 0x400
	return AtomicBit{apb1reg, n}
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

func APB_SetLPEnabled(_ unsafe.Pointer, _ bool) {}
