package usarts

import (
	"unsafe"

	"stm32/usart"
)

var (
	USART1 = (*usart.Dev)(unsafe.Pointer(uintptr(0x40011000)))
	USART2 = (*usart.Dev)(unsafe.Pointer(uintptr(0x40004400)))
	USART3 = (*usart.Dev)(unsafe.Pointer(uintptr(0x40004800)))
	UART4  = (*usart.Dev)(unsafe.Pointer(uintptr(0x40004C00)))
	UART5  = (*usart.Dev)(unsafe.Pointer(uintptr(0x40005000)))
	USART6 = (*usart.Dev)(unsafe.Pointer(uintptr(0x40011400)))
	UART7  = (*usart.Dev)(unsafe.Pointer(uintptr(0x40007800)))
	UART8  = (*usart.Dev)(unsafe.Pointer(uintptr(0x40007C00)))
)
