// +build f030x6

package internal

import (
	"unsafe"

	"stm32/hal/system"

	"stm32/hal/raw/mmap"
)

// Bus returns bus for given peripheral base address.
func Bus(paddr unsafe.Pointer) system.Bus {
	a := uintptr(paddr)
	switch {
	case a >= mmap.AHBPERIPH_BASE:
		return system.AHB
	case a >= mmap.APBPERIPH_BASE:
		return system.APB1
	}
	return -1
}
