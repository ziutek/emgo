package gpio

import (
	"unsafe"
)

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

func (p *Port) Pin(id int) Pin {
	ptr := uintptr(unsafe.Pointer(p))
	return Pin{ptr | uintptr(id&0xf)}
}

// Pin represents one phisical pin (specific pin in specific port).
type Pin struct {
	h uintptr
}

// IsValid reports whether p represents a valid pin.
func (p Pin) IsValid() bool {
	return p.h != 0
}

// Port returns the port where the pin is located.
func (p Pin) Port() *Port {
	return (*Port)(unsafe.Pointer(p.h &^ 0xf))
}

func (p Pin) index() uintptr {
	return p.h & 0xf
}

// Index returns pin index in the port.
func (p Pin) Index() int {
	return int(p.index())
}

// Mask return bitmask that represents the pin.
func (p Pin) Mask() Pins {
	return Pin0 << p.index()
}

// Load returns input value of the pin.
func (p Pin) Load() int {
	return int(p.Port().idr.Load()) >> p.index() & 1
}

// LoadOut returns output value of the pin.
func (p Pin) LoadOut() int {
	return int(p.Port().odr.Load()) >> p.index() & 1
}

// Set sets output value of the pin to 1 in one atomic operation.
func (p Pin) Set() {
	p.Port().bsrr.Store(uint32(Pin0) << p.index())
}

// Clear sets output value of the pin to 0 in one atomic operation.
func (p Pin) Clear() {
	p.Port().bsrr.Store(uint32(Pin0) << 16 << p.index())
}

// Store sets output value of the pin to the least significant bit of val.
func (p Pin) Store(val int) {
	n := p.index()
	v := ^uint32(val)&1<<16 | uint32(val)&1
	p.Port().bsrr.Store(v << n)
}
