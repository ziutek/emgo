package usart

import "unsafe"

type USART struct {
	s   uint32 `C:"volatile"`
	d   uint32 `C:"volatile"`
	br  uint32 `C:"volatile"`
	c1  uint32 `C:"volatile"`
	c2  uint32 `C:"volatile"`
	c3  uint32 `C:"volatile"`
	gtp uint32 `C:"volatile"`
}

var (
	USART1 = (*USART)(unsafe.Pointer(uintptr(0x40011000)))
	USART2 = (*USART)(unsafe.Pointer(uintptr(0x40004400)))
	USART3 = (*USART)(unsafe.Pointer(uintptr(0x40004800)))
	UART4  = (*USART)(unsafe.Pointer(uintptr(0x40004C00)))
	UART5  = (*USART)(unsafe.Pointer(uintptr(0x40005000)))
	USART6 = (*USART)(unsafe.Pointer(uintptr(0x40011400)))
	UART7  = (*USART)(unsafe.Pointer(uintptr(0x40007800)))
	UART8  = (*USART)(unsafe.Pointer(uintptr(0x40007C00)))
)
