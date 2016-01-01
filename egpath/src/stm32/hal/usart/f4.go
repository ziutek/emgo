// +build f40_41xxx f411xe

package usart

import (
	"mmio"
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

var (
	USART1 = (*USART)(unsafe.Pointer(mmap.USART1_BASE))
	USART2 = (*USART)(unsafe.Pointer(mmap.USART2_BASE))
	USART3 = (*USART)(unsafe.Pointer(mmap.USART3_BASE))
	UART4  = (*USART)(unsafe.Pointer(mmap.UART4_BASE))
	UART5  = (*USART)(unsafe.Pointer(mmap.UART5_BASE))
	USART6 = (*USART)(unsafe.Pointer(mmap.USART6_BASE))
	UART7  = (*USART)(unsafe.Pointer(mmap.UART7_BASE))
	UART8  = (*USART)(unsafe.Pointer(mmap.UART8_BASE))
)

func enbit(u *USART) mmio.UM32 {
	var u32 *mmio.U32
	a := u.BaseAddr()
	if a >= mmap.APB2PERIPH_BASE {
		u32 = &rcc.RCC.APB2ENR.U32
		a -= mmap.APB2PERIPH_BASE
	} else {
		u32 = &rcc.RCC.APB1ENR.U32
		a -= mmap.APB1PERIPH_BASE
	}
	return mmio.UM32{u32, 1 << (a / 0x400)}
}

func lpenbit(u *USART) mmio.UM32 {
	var u32 *mmio.U32
	a := u.BaseAddr()
	if a >= mmap.APB2PERIPH_BASE {
		u32 = &rcc.RCC.APB2LPENR.U32
		a -= mmap.APB2PERIPH_BASE
	} else {
		u32 = &rcc.RCC.APB1LPENR.U32
		a -= mmap.APB1PERIPH_BASE
	}
	return mmio.UM32{u32, 1 << (a / 0x400)}
}

func rstbit(u *USART) mmio.UM32 {
	var u32 *mmio.U32
	a := u.BaseAddr()
	if a >= mmap.APB2PERIPH_BASE {
		u32 = &rcc.RCC.APB2RSTR.U32
		a -= mmap.APB2PERIPH_BASE
	} else {
		u32 = &rcc.RCC.APB1RSTR.U32
		a -= mmap.APB1PERIPH_BASE
	}
	return mmio.UM32{u32, 1 << (a / 0x400)}
}
