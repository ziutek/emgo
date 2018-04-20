// +build f030x6

package spi

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	SPI1 = (*Periph)(unsafe.Pointer(mmap.SPI1_BASE))
)
