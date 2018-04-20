// +build f030x8

package spi

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	SPI1 = (*Periph)(unsafe.Pointer(mmap.SPI1_BASE))
	SPI2 = (*Periph)(unsafe.Pointer(mmap.SPI2_BASE))
)
