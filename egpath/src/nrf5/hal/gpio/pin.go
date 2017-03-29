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
	return (*Port)(unsafe.Pointer(p.h &^ 0x1f))
}

func (p Pin) index() uintptr {
	return p.h & 0x1f
}

// Index returns pin index in the port.
func (p Pin) Index() int {
	return int(p.index())
}

// Setup configures pin.
func (p Pin) Setup(cfg Config) {
	p.Port().SetupPin(p.Index(), cfg)
}

// Config returns current configuration of pin.
func (p Pin) Config() Config {
	return p.Port().PinConfig(p.Index())
}

// Mask returns bitmask that represents the pin.
func (p Pin) Mask() Pins {
	return Pin0 << p.index()
}

// Load returns input value of the pin.
func (p Pin) Load() int {
	return int(p.Port().in.Load()) >> p.index() & 1
}

// LoadOut returns output value of the pin.
func (p Pin) LoadOut() int {
	return int(p.Port().out.Load()) >> p.index() & 1
}

// Set sets output value of the pin to 1 in one atomic operation.
func (p Pin) Set() {
	p.Port().outset.Store(uint32(Pin0) << p.index())
}

// Clear sets output value of the pin to 0 in one atomic operation.
func (p Pin) Clear() {
	p.Port().outclr.Store(uint32(Pin0) << p.index())
}

// Store sets output value of the pin to the least significant bit of val.
func (p Pin) Store(val int) {
	port := p.Port()
	n := p.index()
	if val&1 != 0 {
		port.outset.Store(uint32(Pin0) << n)
	} else {
		port.outclr.Store(uint32(Pin0) << n)
	}
}
