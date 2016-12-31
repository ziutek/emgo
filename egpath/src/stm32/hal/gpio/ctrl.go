package gpio

// Pins is a bitmask which represents the pins of GPIO port.
type Pins uint16

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
)

// Pins returns input value of pins.
func (p *Port) Pins(pins Pins) Pins {
	return Pins(p.idr.Bits(uint16(pins)))
}

// PinsOut returns output value of pins.
func (p *Port) PinsOut(pins Pins) Pins {
	return Pins(p.odr.Bits(uint16(pins)))
}

// Set sets output value of pins to 1 in one atomic operation.
func (p *Port) SetPins(pins Pins) {
	p.bsrr.Store(uint32(pins))
}

// Clear sets output value of pins to 0 in one atomic operation.
func (p *Port) ClearPins(pins Pins) {
	p.bsrr.Store(uint32(pins) << 16)
}

// ClearAndSet clears and sets output value of all pins in one atomic operation.
// Setting pins has priority above clearing bits.
func (p *Port) ClearAndSet(clear, set Pins) {
	p.bsrr.Store(uint32(clear)<<16 | uint32(set))
}

// StorePins sets pins specified by pins to val in one atomic operation.
func (p *Port) StorePins(pins, val Pins) {
	m := uint32(pins)<<16 | uint32(pins)
	v := ^uint32(val)<<16 | uint32(val)
	p.bsrr.Store(m & v)
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
	p.odr.Store(uint16(val))
}
