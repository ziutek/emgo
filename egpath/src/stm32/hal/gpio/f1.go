// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package gpio

import (
	"mmio"
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
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
)

type registers struct {
	cr   [2]mmio.U32
	idr  mmio.U16
	_    uint16
	odr  mmio.U16
	_    uint16
	bsrr mmio.U32
	brr  mmio.U32
	lckr mmio.U32
}

const (
	// Mode
	out   = 1
	alt   = 8 + 1
	altIn = 0
	ana   = 2

	// Pull
	pullUp   = 8 + 4 + 1
	pullDown = 8 + 4 + 0

	// Driver
	openDrain = 4

	// Speed
	veryLow  = 1
	low      = 1
	high     = 2
	veryHigh = 2
)

func enableClock(p *Port, _ bool) {
	rcc.RCC.APB2ENR.U32.AtomicSetBits(uint32(rcc.IOPAEN << portnum(p)))
	rcc.RCC.APB2ENR.Load() // RCC delay (workaround for silicon bugs).
}

func disableClock(p *Port) {
	rcc.RCC.APB2ENR.U32.AtomicClearBits(uint32(rcc.IOPAEN << portnum(p)))
}

func reset(p *Port) {
	pnum := portnum(p)
	rcc.RCC.APB2RSTR.U32.AtomicSetBits(uint32(rcc.IOPARST << pnum))
	rcc.RCC.APB2RSTR.U32.AtomicClearBits(uint32(rcc.IOPARST << pnum))
}

func setup(p *Port, n int, cfg *Config) {
	cr := &p.cr[n>>3]
	pos := uint(n & 7 * 4)
	sel := uint32(0xf) << pos
	switch {
	case cfg.Mode == 0: // In, AltIn.
		cm := uint32(cfg.Pull)&(8+4) ^ 4
		cr.StoreBits(sel, cm<<pos)
		p.StorePins(Pin0<<uint(n), Pins(cfg.Pull)<<uint(n))
	case cfg.Mode&1 != 0: // Out, Alt.
		cm := uint32(cfg.Mode) & 8
		cm |= uint32(cfg.Driver)
		cm |= uint32(cfg.Speed) + 1
		cr.StoreBits(sel, cm<<pos)
	default: // Ana.
		cr.ClearBits(sel)
	}
}
