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

type Port struct {
	cr   [2]mmio.U32
	idr  mmio.U32
	odr  mmio.U32
	bsrr mmio.U32
	brr  mmio.U32
	lckr mmio.U32
}

// Num returns ordinar number of p on its peripheral bus.
func (p *Port) Num() int {
	return pnum(p)
}

// EnableClock enables clock for port p.
func (p *Port) EnableClock(_ bool) {
	enr().SetBit(pnum(p))
	_ = enr().Load() // Workaround (RCC delay).
}

// DisableClock disables clock for port p.
func (p *Port) DisableClock() {
	enr().ClearBit(pnum(p))
}

// Reset resets port p.
func (p *Port) Reset() {
	n := pnum(p)
	rstr().SetBit(n)
	rstr().ClearBit(n)
}

type Mode byte

const (
	In      Mode = 0
	Out     Mode = 1
	OutSlow Mode = 2
	OutFast Mode = 3
)

// SetMode sets I/O mode for n-th pin.
func (p *Port) SetMode(n int, mode Mode) {
	cr := &p.cr[n>>3]
	m := uint(n & 7 * 4)
	cr.StoreBits(0xf<<m, uint32(mode)<<m)
}

// PinIn returns current input value of n-th pin.
func (p *Port) PinIn(n int) int {
	return p.idr.Bit(n)
}

// PinOut returns current output value of n-th pin.
func (p *Port) PinOut(n int) int {
	return p.odr.Bit(n)
}

// SetPin sets output value of n-th pin to 1.
func (p *Port) SetPin(n int) {
	p.bsrr.Store(1 << uint(n))
}

// ClearPin sets output value of n-th pin to 0.
func (p *Port) ClearPin(n int) {
	p.bsrr.Store(0x10000 << uint(n))
}

// SetPins sets output value of pins on positions specified by mask to 1.
func (p *Port) SetPins(mask uint16) {
	p.bsrr.Store(uint32(mask))
}

// ClearPins sets output value of pins on positions specified by mask to 0.
func (p *Port) ClearPins(mask uint16) {
	p.bsrr.Store(uint32(mask) << 16)
}

// ClearAndSet clears and sets output value of all pins on positions specified/
// by csmask.
// High 16 bits in csmask specifies which pins should be 0.
// Low 16 bits in csmask specifies which pins should be 1.
// Setting bits in csmask has priority above clearing bits.
func (p *Port) ClearAndSet(csmask uint32) {
	p.bsrr.Store(csmask)
}

// LoadIn returns input value of pins.
func (p *Port) LoadIn() uint16 {
	return uint16(p.idr.Load())
}

// LoadOut returns output value of pins.
func (p *Port) LoadOut() uint16 {
	return uint16(p.odr.Load())
}

// Store sets output value of all pins to value specified by bits.
func (p *Port) Store(bits uint16) {
	p.odr.Store(uint32(bits))
}

func pnum(p *Port) int {
	return int(uintptr(unsafe.Pointer(p))-mmap.APB2PERIPH_BASE) / 0x400
}

func enr() *mmio.U32  { return &rcc.RCC.APB2ENR.U32 }
func rstr() *mmio.U32 { return &rcc.RCC.APB2RSTR.U32 }
