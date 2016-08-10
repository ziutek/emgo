// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package i2c

import (
	"unsafe"
	
	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	I2C1 = (*Periph)(unsafe.Pointer(mmap.I2C1_BASE))
	I2C2 = (*Periph)(unsafe.Pointer(mmap.I2C2_BASE))
)
