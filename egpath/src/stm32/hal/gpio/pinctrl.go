package gpio

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

// Set sets output value of pins to 1.
func (p *Port) Set(pins Pins) {
	p.bsrr.Store(uint32(pins))
}

// Clear sets output value of pins to 0.
func (p *Port) Clear(pins Pins) {
	p.bsrr.Store(uint32(pins) << 16)
}

// ClearAndSet clears and sets output value of all pins on positions specified
// by cspins. Upper half of cspins specifies which pins should be 0. Lower half
// of cspins specifies which pins should be 1. Setting bits in cspins has
// priority above clearing bits.
func (p *Port) ClearAndSet(cspins Pins) {
	p.bsrr.Store(uint32(cspins))
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

// Mask returns input value of pins specified by mask.
func (p *Port) Mask(mask Pins) Pins {
	return Pins(p.idr.Bits(uint32(mask)))
}

// MaskOut returns output value of pins specified by mask.
func (p *Port) MaskOut(mask Pins) Pins {
	return Pins(p.odr.Bits(uint32(mask)))
}

// MaskStore sets pins specified by mask to val.
func (p *Port) MaskStore(mask, val Pins) {
	mask |= mask << 16
	val |= ^val << 16
	p.bsrr.Store(uint32(mask & val))
}

/*
// Pin returns current input value of n-th pin.
func (p *Port) Pin(n int) int {
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

// StorePin sets output value of n-th pin to v&1.
func (p *Port) StorePin(n, v int) {
	n += (1 - v&1) * 16
	p.bsrr.Store(1 << uint(n))
}
*/
