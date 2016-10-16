// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl l1xx_md l1xx_mdp l1xx_hd l1xx_xl

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
)
