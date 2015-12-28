// +build l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package gpio

import (
	"mmio"
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

const (
	VeryLow Speed = 0 // 400 kHz (CL = 50 pF, VDD > 2.7 V)
	Low     Speed = 1 //   2 MHz (CL = 50 pF, VDD > 2.7 V)
	Medium  Speed = 2 //  10 MHz (CL = 50 pF, VDD > 2.7 V)
	High    Speed = 3 //  50 MHz (CL = 50 pF, VDD > 2.7 V)
)

var (
	A = (*Port)(unsafe.Pointer(mmap.GPIOA_BASE))
	B = (*Port)(unsafe.Pointer(mmap.GPIOB_BASE))
	C = (*Port)(unsafe.Pointer(mmap.GPIOC_BASE))
	D = (*Port)(unsafe.Pointer(mmap.GPIOD_BASE))
	E = (*Port)(unsafe.Pointer(mmap.GPIOE_BASE))
	F = (*Port)(unsafe.Pointer(mmap.GPIOF_BASE))
	G = (*Port)(unsafe.Pointer(mmap.GPIOG_BASE))
)

func pnum(p *Port) int {
	return int(uintptr(unsafe.Pointer(p))-mmap.AHBPERIPH_BASE) / 0x400
}

func enr() *mmio.U32   { return (*mmio.U32)(&rcc.RCC.AHBENR.U32) }
func lpenr() *mmio.U32 { return (*mmio.U32)(&rcc.RCC.AHBLPENR.U32) }
func rstr() *mmio.U32  { return (*mmio.U32)(&rcc.RCC.AHBRSTR.U32) }
