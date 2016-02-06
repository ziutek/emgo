package gpio

import (
	"arch/cortexm/bitband"
)

// Pins is bitmask which lower 16 bis represents pins of GPIO port.
type Pins uint32

const (
	Pin0 Pins = 1 << iota
	Pin1
	Pin2
	Pin3
	Pin4
	Pin5
	Pin6
	Pin7
	Pin8
	Pin9
	Pin10
	Pin11
	Pin12
	Pin13
	Pin14
	Pin15
	All Pins = 0xffff
)

// Pins returns input value of pins.
func (p *Port) Pins(pins Pins) Pins {
	return Pins(p.idr.Bits(uint32(pins)))
}

// PinsOut returns output value of pins.
func (p *Port) PinsOut(pins Pins) Pins {
	return Pins(p.odr.Bits(uint32(pins)))
}

// Set sets output value of pins to 1.
func (p *Port) SetPins(pins Pins) {
	p.bsrr.Store(uint32(pins))
}

// Clear sets output value of pins to 0.
func (p *Port) ClearPins(pins Pins) {
	p.bsrr.Store(uint32(pins) << 16)
}

// ClearAndSet clears and sets output value of all pins on positions specified
// by cspins. Upper half of cspins specifies which pins should be 0. Lower half
// of cspins specifies which pins should be 1. Setting bits in cspins has
// priority above clearing bits.
func (p *Port) ClearAndSet(cspins Pins) {
	p.bsrr.Store(uint32(cspins))
}

// StorePins sets pins specified by pins to val.
func (p *Port) StorePins(pins, val Pins) {
	pins |= pins << 16
	val |= ^val << 16
	p.bsrr.Store(uint32(pins & val))
}

// Load returns input value of all pins.
func (p *Port) Load() Pins {
	return Pins(p.idr.Load())
}

// LoadOut returns output value of all pins.
func (p *Port) LoadOut() Pins {
	return Pins(p.odr.Load())
}

// Store sets output value of all pins to value specified by val.
func (p *Port) Store(val Pins) {
	p.odr.Store(uint32(val))
}

// Pin returns bitband alias to input value of n-th pin.
func (p *Port) InPin(n int) bitband.Bit {
	return bitband.Alias32(&p.idr).Bit(n)
}

// OutPin returns bitband alias to output value of n-th pin.
func (p *Port) OutPin(n int) bitband.Bit {
	return bitband.Alias32(&p.odr).Bit(n)
}
