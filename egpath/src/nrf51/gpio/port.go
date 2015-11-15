package gpio

import (
	"unsafe"

	"nrf51/periph"
)

type Port struct {
	out    uint32
	outset uint32
	outclr uint32
	in     uint32
	dir    uint32
	dirset uint32
	dirclr uint32
	_      [120]uint32
	pincnf [32]uint32
} //c:volatile

var P0 = (*Port)(unsafe.Pointer(periph.BaseAHB + 0x504))

type Mode byte

const (
	In Mode = iota
	InOut
	Discon
	Out
)

// Mode returns I/O mode for n-th pin.
func (p *Port) Mode(n int) Mode {
	return Mode(p.pincnf[n] & 3)
}

// SetMode sets I/O mode for n-th pin.
func (p *Port) SetMode(n int, mode Mode) {
	p.pincnf[n] = p.pincnf[n]&^3 | uint32(mode)
}

// InPin returns current value of n-th input pin.
func (p *Port) InPin(n int) int {
	return int(p.in>>uint(n)) & 1
}

// OutPin returns current value of n-th output pin.
func (p *Port) OutPin(n int) int {
	return int(p.out>>uint(n)) & 1
}

// SetPin sets n-th output pin to 1.
func (p *Port) SetPin(n int) {
	p.outset = uint32(1) << uint(n)
}

// ClearOutPin sets n-th output pin to 0.
func (p *Port) ClearPin(n int) {
	p.outclr = uint32(1) << uint(n)
}

// SetPins sets output pins on positions specified by pins to 1.
func (p *Port) SetPins(bits uint32) {
	p.outset = bits
}

// ClearPins sets output pins on positions specified by bits to 0.
func (p *Port) ClearPins(bits uint32) {
	p.outclr = bits
}

// LoadIn returns value of input pins.
func (p *Port) LoadIn() uint32 {
	return p.in
}

// LoadOut returns value of output pins.
func (p *Port) LoadOut() uint32 {
	return p.out
}

// Store sets output pins to value specified by bits.
func (p *Port) Store(bits uint32) {
	p.out = bits
}
