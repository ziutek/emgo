// +build f40_41xxx f411xe

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
	USART6 = (*USART)(unsafe.Pointer(mmap.USART6_BASE))
	UART7  = (*USART)(unsafe.Pointer(mmap.UART7_BASE))
	UART8  = (*USART)(unsafe.Pointer(mmap.UART8_BASE))
)
