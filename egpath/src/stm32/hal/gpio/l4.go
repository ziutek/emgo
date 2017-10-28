// +build l476xx

package gpio

import (
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

const ()

const (
	veryLow  = -1 // Not supported.
	low      = -1 //   v MHz (CL = 50 pF, VDD > 2.7 V)
	_        = 0  //  x MHz (CL = 50 pF, VDD > 2.7 V)
	high     = 1  //  y MHz (CL = 40 pF, VDD > 2.7 V)
	veryHigh = 2  // z MHz (CL = 30 pF, VDD > 2.7 V)
)

//emgo:const
var (
	A = (*Port)(unsafe.Pointer(mmap.GPIOA_BASE))
	B = (*Port)(unsafe.Pointer(mmap.GPIOB_BASE))
	C = (*Port)(unsafe.Pointer(mmap.GPIOC_BASE))
	D = (*Port)(unsafe.Pointer(mmap.GPIOD_BASE))
	E = (*Port)(unsafe.Pointer(mmap.GPIOE_BASE))
	F = (*Port)(unsafe.Pointer(mmap.GPIOF_BASE))
	G = (*Port)(unsafe.Pointer(mmap.GPIOG_BASE))
	H = (*Port)(unsafe.Pointer(mmap.GPIOH_BASE))
)

func enreg() *rcc.AHB2ENR   { return &rcc.RCC.AHB2ENR }
func rstreg() *rcc.AHB2RSTR { return &rcc.RCC.AHB2RSTR }

func lpenaclk(pnum uint) {
	rcc.RCC.AHB2SMENR.U32.AtomicSetBits(uint32(rcc.GPIOASMEN << pnum))
}
func lpdisclk(pnum uint) {
	rcc.RCC.AHB2SMENR.U32.AtomicClearBits(uint32(rcc.GPIOASMEN << pnum))
}
