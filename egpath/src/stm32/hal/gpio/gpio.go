// +build f40_41xxx f411xe l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package gpio

import (
	"mmio"
)

// Port represents GPIO port.
type Port struct {
	moder   mmio.U32
	otyper  mmio.U32
	ospeedr mmio.U32
	pupdr   mmio.U32
	idr     mmio.U32
	odr     mmio.U32
	bsrr    mmio.U32
	lckr    mmio.U32
	afr     [2]mmio.U32
}

// Enable enables clock for port p.
// lp determines whether the clock remains on in low power (sleep) mode.
func (p *Port) Enable(lp bool) {
	p.enable(lp)
}

// Disable disables clock for port p.
func (p *Port) Disable() {
	p.disable()
}

// Reset resets port p.
func (p *Port) Reset() {
	p.reset()
}

// Mode represents pin operation mode (input, output, alternate, analog).
type Mode byte

const (
	In  Mode = 0 // Input mode.
	Out Mode = 1 // Output mode.
	Alt Mode = 2 // Alternate function mode.
	Ana Mode = 3 // Analog mode.
)

// Mode returns mode for n-th pin.
func (p *Port) Mode(n int) Mode {
	return Mode(p.moder.Load()>>uint(n*2)) & 3
}

// SetMode sets I/O mode for n-th pin.
func (p *Port) SetMode(n int, mode Mode) {
	m := uint(n * 2)
	p.moder.StoreBits(3<<m, uint32(mode)<<m)
}

// OutType represents type of output pin.
type OutType byte

const (
	PushPull  OutType = 0
	OpenDrain OutType = 1
)

// OutType returns current type of n-th pin.
func (p *Port) OutType(n int) OutType {
	return OutType(p.otyper.Bit(n))
}

// SetOutType sets type for n-th output pin.
func (p *Port) SetOutType(n int, otyp OutType) {
	p.otyper.StoreBit(n, int(otyp))
}

// Speed
type Speed byte

// OutSpeed return current speed for n-th pin.
func (p *Port) OutSpeed(n int) Speed {
	return Speed(p.ospeedr.Load()>>uint(n*2)) & 3
}

// SetOutSpeed sets speed for n-th output pin.
func (p *Port) SetOutSpeed(n int, speed Speed) {
	m := uint(n * 2)
	p.ospeedr.StoreBits(3<<m, uint32(speed)<<m)
}

type Pull byte

const (
	Float    Pull = 0
	PullUp   Pull = 1
	PullDown Pull = 2
)

// Pull returns current pull configuration of n-th pin.
func (p *Port) Pull(n int) Pull {
	return Pull(p.pupdr.Load()>>uint(n*2)) & 3
}

// SetPull sets internal pull-up/pull-down cirquits for n-th output pin.
func (p *Port) SetPull(n int, pull Pull) {
	m := uint(n * 2)
	p.pupdr.StoreBits(3<<m, uint32(pull)<<m)
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

/*
// AltFunc returns current alternate function for n-th pin in port g.
func (p *Port) AltFunc(n int) AltFunc {
	m := uint(n & 7 * 4)
	n = n >> 3 & 1
	return AltFunc(p.afr[n].Load() >> m & 0xf)
}

// SetAltFunc sets alternate function af for n-th pin in port g.
func (p *Port) SetAltFunc(n int, af AltFunc) {
	m := uint(n & 7 * 4)
	n = n >> 3 & 1
	p.afr[n].StoreBits(0xf<<m, uint32(af)<<m)
}
*/
