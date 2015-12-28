// +build f40_41xxx f411xe

package gpio

import (
	"mmio"
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

const (
	Low      Speed = 0 //   2 MHz (CL = 50 pF, VDD > 2.7 V)
	Medium   Speed = 1 //  25 MHz (CL = 50 pF, VDD > 2.7 V)
	High     Speed = 2 //  50 MHz (CL = 40 pF, VDD > 2.7 V)
	VeryHigh Speed = 3 // 100 MHz (CL = 30 pF, VDD > 2.7 V)
)

var (
	A = (*Port)(unsafe.Pointer(mmap.GPIOA_BASE))
	B = (*Port)(unsafe.Pointer(mmap.GPIOB_BASE))
	C = (*Port)(unsafe.Pointer(mmap.GPIOC_BASE))
	D = (*Port)(unsafe.Pointer(mmap.GPIOD_BASE))
	E = (*Port)(unsafe.Pointer(mmap.GPIOE_BASE))
	F = (*Port)(unsafe.Pointer(mmap.GPIOF_BASE))
	G = (*Port)(unsafe.Pointer(mmap.GPIOG_BASE))
	H = (*Port)(unsafe.Pointer(mmap.GPIOH_BASE))
	I = (*Port)(unsafe.Pointer(mmap.GPIOI_BASE))
	J = (*Port)(unsafe.Pointer(mmap.GPIOJ_BASE))
	K = (*Port)(unsafe.Pointer(mmap.GPIOK_BASE))
)

func pnum(p *Port) int {
	return int(uintptr(unsafe.Pointer(p))-mmap.AHB1PERIPH_BASE) / 0x400
}

func enr() *mmio.U32   { return (*mmio.U32)(&rcc.RCC.AHB1ENR.U32) }
func lpenr() *mmio.U32 { return (*mmio.U32)(&rcc.RCC.AHB1LPENR.U32) }
func rstr() *mmio.U32  { return (*mmio.U32)(&rcc.RCC.AHB1RSTR.U32) }
