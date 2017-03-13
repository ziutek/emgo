package gpio

import (
	"unsafe"

	"stm32/hal/internal"
	"stm32/hal/system"
)

// Port represents GPIO port.
type Port struct {
	registers
}

// Bus returns a bus to which p is connected.
func (p *Port) Bus() system.Bus {
	return internal.Bus(unsafe.Pointer(p))
}

// Num returns port number: A.Num() == 0, B.Num() = 1, ...
func (p *Port) Num() int {
	return int(portnum(p))
}

// EnableClock enables clock for port p.
// lp determines whether the clock remains on in low power (sleep) mode.
func (p *Port) EnableClock(lp bool) {
	enableClock(p, lp)
}

// DisableClock disables clock for port p.
func (p *Port) DisableClock() {
	disableClock(p)
}

// Reset resets port p.
func (p *Port) Reset() {
	reset(p)
}

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

func (p *Port) Pin(id int) Pin {
	ptr := uintptr(unsafe.Pointer(p))
	return Pin{ptr | uintptr(id&0xf)}
}
