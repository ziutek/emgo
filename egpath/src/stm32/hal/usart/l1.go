// +build l1xx_md l1xx_mdp l1xx_hd l1xx_xl

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
