// +build f746xx l476xx

package sdmmc

import (
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/sdmmc"
)

type Periph struct {
	raw sdmmc.SDMMC_Periph
}

//emgo:const
var SDMMC = (*Periph)(unsafe.Pointer(mmap.SDMMC_BASE))
