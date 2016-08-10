// +build f40_41xxx f411xe

package usart

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	USART1 = (*Periph)(unsafe.Pointer(mmap.USART1_BASE))
	USART2 = (*Periph)(unsafe.Pointer(mmap.USART2_BASE))
	USART3 = (*Periph)(unsafe.Pointer(mmap.USART3_BASE))
	UART4  = (*Periph)(unsafe.Pointer(mmap.UART4_BASE))
	UART5  = (*Periph)(unsafe.Pointer(mmap.UART5_BASE))
	USART6 = (*Periph)(unsafe.Pointer(mmap.USART6_BASE))
	UART7  = (*Periph)(unsafe.Pointer(mmap.UART7_BASE))
	UART8  = (*Periph)(unsafe.Pointer(mmap.UART8_BASE))
)
