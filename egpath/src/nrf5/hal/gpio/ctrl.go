package gpio

// Pins returns input value of pins.
func (p *Port) Pins(pins Pins) Pins {
	return Pins(p.in.Bits(uint32(pins)))
}

// PinsOut returns output value of pins.
func (p *Port) PinsOut(pins Pins) Pins {
	return Pins(p.out.Bits(uint32(pins)))
}

// SetPins sets output value of pins to 1 in one atomic operation.
func (p *Port) SetPins(pins Pins) {
	p.outset.Store(uint32(pins))
}

// ClearPins sets output value of pins to 0 in one atomic operation.
func (p *Port) ClearPins(pins Pins) {
	p.outclr.Store(uint32(pins))
}

// Load returns input value of all pins.
func (p *Port) Load() Pins {
	return Pins(p.in.Load())
}

// LoadOut returns output value of all pins.
func (p *Port) LoadOut() Pins {
	return Pins(p.out.Load())
}

// Store sets output value of all pins to value specified by val.
func (p *Port) Store(val Pins) {
	p.out.Store(uint32(val))
}

/*
// InPin returns current value of n-th input pin.
func (p *Port) InPin(n int) int {
	return int(p.in.Load()>>uint(n)) & 1
}

// OutPin returns current value of n-th output pin.
func (p *Port) OutPin(n int) int {
	return int(p.out.Load()>>uint(n)) & 1
}

// SetPin sets output value of n-th pin to 1.
func (p *Port) SetPin(n int) {
	p.outset.Store(1 << uint(n))
}

// ClearPin sets output value of n-th pin to 0.
func (p *Port) ClearPin(n int) {
	p.outclr.Store(1 << uint(n))
}
*/
