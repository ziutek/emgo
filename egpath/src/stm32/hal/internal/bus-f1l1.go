// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl  l1xx_md l1xx_mdp l1xx_hd l1xx_xl

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
	case a >= mmap.APB2PERIPH_BASE:
		return system.APB2
	case a >= mmap.APB1PERIPH_BASE:
		return system.APB1
	}
	return -1
}
