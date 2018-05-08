// +build f030x6 f030x8

package gpio

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	A = (*Port)(unsafe.Pointer(mmap.GPIOA_BASE))
	B = (*Port)(unsafe.Pointer(mmap.GPIOB_BASE))
	C = (*Port)(unsafe.Pointer(mmap.GPIOC_BASE))
	D = (*Port)(unsafe.Pointer(mmap.GPIOD_BASE))
	F = (*Port)(unsafe.Pointer(mmap.GPIOF_BASE))
)
