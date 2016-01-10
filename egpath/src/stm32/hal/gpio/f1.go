// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package gpio

import (
	"mmio"
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
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

type registers struct {
	cr   [2]mmio.U32
	idr  mmio.U32
	odr  mmio.U32
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

func pnum(p *Port) int {
	return int(uintptr(unsafe.Pointer(p))-mmap.APB2PERIPH_BASE) / 0x400
}

func enr() *mmio.U32  { return &rcc.RCC.APB2ENR.U32 }
func rstr() *mmio.U32 { return &rcc.RCC.APB2RSTR.U32 }

func enableClock(p *Port, _ bool) {
	enr().SetBit(pnum(p))
	_ = enr().Load() // RCC delay (workaround for silicon bugs).
}

func setup(p *Port, n int, cfg *Config) {
	cr := &p.cr[n>>3]
	pos := uint(n & 7 * 4)
	sel := uint32(0xf) << pos
	switch {
	case cfg.Mode == 0: // In, AltIn.
		cm := uint32(cfg.Pull)&(8+4) ^ 4
		cr.StoreBits(sel, cm<<pos)
		p.MaskStore(Pin0<<uint(n), Pins(cfg.Pull)<<uint(n))
	case cfg.Mode&1 != 0: // Out, Alt.
		cm := uint32(cfg.Mode) & 8
		cm |= uint32(cfg.Driver)
		cm |= uint32(cfg.Speed) + 1
		cr.StoreBits(sel, cm<<pos)
	default: // Ana.
		cr.ClearBits(sel)
	}
}

/*
func setModeIn(p *Port, n int, m *ModeIn) {
	cr := &p.cr[n>>3]
	pos := uint(n & 7 * 4)
	cnf := uint32(m.Pull) & 8
	cr.StoreBits(0xf<<pos, cnf<<pos)
	p.StorePin(n, int(m.Pull))
	// Ignore mo.Alt.
}

func setModeOut(p *Port, n int, m *ModeOut) {
	cr := &p.cr[n>>3]
	pos := uint(n & 7 * 4)
	cm := uint32(bits.One(m.Alt)) << 3
	cm |= uint32(m.Driver)
	cm |= uint32(m.Speed) + 1
	cr.StoreBits(0xf<<pos, cm<<pos)
	// Ignore mo.Pull (not supported by STM32F1xx).
}

func setModeAna(p *Port, n int) {
	cr := &p.cr[n>>3]
	pos := uint(n & 7 * 4)
	cr.ClearBits(0xf << pos)
}
*/
