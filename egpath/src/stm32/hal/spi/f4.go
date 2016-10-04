// +build f40_41xxx f411xe

package spi

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	SPI1 = (*Periph)(unsafe.Pointer(mmap.SPI1_BASE))
	SPI2 = (*Periph)(unsafe.Pointer(mmap.SPI2_BASE))
	SPI3 = (*Periph)(unsafe.Pointer(mmap.SPI3_BASE))
	SPI4 = (*Periph)(unsafe.Pointer(mmap.SPI4_BASE))
)
