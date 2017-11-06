// +build l476xx

package usart

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	USART1  = (*Periph)(unsafe.Pointer(mmap.USART1_BASE))
	USART2  = (*Periph)(unsafe.Pointer(mmap.USART2_BASE))
	USART3  = (*Periph)(unsafe.Pointer(mmap.USART3_BASE))
	UART4   = (*Periph)(unsafe.Pointer(mmap.UART4_BASE))
	UART5   = (*Periph)(unsafe.Pointer(mmap.UART5_BASE))
	LPUART1 = (*Periph)(unsafe.Pointer(mmap.LPUART1_BASE))
)
