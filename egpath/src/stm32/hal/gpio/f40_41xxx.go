// +build f40_41xxx

package gpio

import (
	"unsafe"

	"stm32/o/f40_41xxx/mmap"
	"stm32/o/f40_41xxx/rcc"
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

func (p *Port) n() int {
	return int(uintptr(unsafe.Pointer(p))-mmap.AHB1PERIPH_BASE) / 0x400
}

func (p *Port) enable(lp bool) {
	n := p.n()
	rcc.RCC.AHB1ENR.SetBit(n)
	if lp {
		rcc.RCC.AHB1LPENR.SetBit(n)
	} else {
		rcc.RCC.AHB1LPENR.ClearBit(n)
	}
}

func (p *Port) disable() {
	rcc.RCC.AHB1ENR.ClearBit(p.n())
}

func (p *Port) reset() {
	n := p.n()
	rcc.RCC.AHB1RSTR.SetBit(n)
	rcc.RCC.AHB1RSTR.ClearBit(n)
}
