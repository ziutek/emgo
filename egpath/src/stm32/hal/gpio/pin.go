package gpio

import (
	"unsafe"
)

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

// Setup configures pin.
func (p Pin) Setup(cfg *Config) {
	p.Port().SetupPin(p.Index(), cfg)
}

// Lock locks configuration of pin. Locked configuration can not be modified until
// reset
func (p Pin) Lock() {
	p.Port().Lock(Pin0 << p.index())
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
