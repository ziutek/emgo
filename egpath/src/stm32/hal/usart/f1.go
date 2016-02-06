// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package usart

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

var (
	USART1 = (*USART)(unsafe.Pointer(mmap.USART1_BASE))
	USART2 = (*USART)(unsafe.Pointer(mmap.USART2_BASE))
	USART3 = (*USART)(unsafe.Pointer(mmap.USART3_BASE))
	UART4  = (*USART)(unsafe.Pointer(mmap.UART4_BASE))
	UART5  = (*USART)(unsafe.Pointer(mmap.UART5_BASE))
)
