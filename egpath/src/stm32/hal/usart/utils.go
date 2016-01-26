package usart

import (
	"mmio"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

func setEnabled(u *USART, en bool) {
	var u32 *mmio.U32
	a := u.BaseAddr()
	if a >= mmap.APB2PERIPH_BASE {
		u32 = &rcc.RCC.APB2ENR.U32
		a -= mmap.APB2PERIPH_BASE
	} else {
		u32 = &rcc.RCC.APB1ENR.U32
		a -= mmap.APB1PERIPH_BASE
	}
	bit := int(a / 0x400)
	if en {
		u32.SetBit(bit)
	} else {
		u32.ClearBit(bit)
	}
	_ = u32.Load() // Workaround (RCC delay).
}

func reset(u *USART) {
	var u32 *mmio.U32
	a := u.BaseAddr()
	if a >= mmap.APB2PERIPH_BASE {
		u32 = &rcc.RCC.APB2RSTR.U32
		a -= mmap.APB2PERIPH_BASE
	} else {
		u32 = &rcc.RCC.APB1RSTR.U32
		a -= mmap.APB1PERIPH_BASE
	}
	bit := int(a / 0x400)
	u32.SetBit(bit)
	u32.ClearBit(bit)
}
