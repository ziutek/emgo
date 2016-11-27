package gpio

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

func portnum(p *Port) uint {
	return uint(uintptr(unsafe.Pointer(p))-mmap.GPIOA_BASE) / 0x400
}
