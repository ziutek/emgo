// +build f411xe

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
	E = (*Port)(unsafe.Pointer(mmap.GPIOE_BASE))
	H = (*Port)(unsafe.Pointer(mmap.GPIOH_BASE))
)
