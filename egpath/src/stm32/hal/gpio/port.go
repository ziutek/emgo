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
	return portnum(p)
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
