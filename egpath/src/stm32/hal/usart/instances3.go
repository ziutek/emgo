// +build f411xe

package usart

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	USART1 = (*Periph)(unsafe.Pointer(mmap.USART1_BASE))
	USART2 = (*Periph)(unsafe.Pointer(mmap.USART2_BASE))
	USART6 = (*Periph)(unsafe.Pointer(mmap.USART6_BASE))
)
