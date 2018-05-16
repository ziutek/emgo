// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl f303xe f40_41xxx f411xe f746xx l1xx_md l1xx_mdp l1xx_hd l1xx_xl l476xx

package dma

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	DMA1 = (*DMA)(unsafe.Pointer(mmap.DMA1_BASE))
	DMA2 = (*DMA)(unsafe.Pointer(mmap.DMA2_BASE))
)
