// +build f030x6 f030x8

package dma

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	DMA1 = (*DMA)(unsafe.Pointer(mmap.DMA1_BASE))
)
