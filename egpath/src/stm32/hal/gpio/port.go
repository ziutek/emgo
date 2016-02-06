package gpio

import (
	"mmio"
	"unsafe"

	"arch/cortexm/bitband"

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

// Num returns number of p on its peripheral bus.
func (p *Port) Num() int {
	return pnum(p)
}

// EnableClock enables clock for port p.
// lp determines whether the clock remains on in low power (sleep) mode.
func (p *Port) EnableClock(lp bool) {
	enableClock(p, lp)
}

// DisableClock disables clock for port p.
func (p *Port) DisableClock() {
	bit(p, enreg()).Clear()
}

// Reset resets port p.
func (p *Port) Reset() {
	bit := bit(p, rstreg())
	bit.Set()
	bit.Clear()
}

func bit(p *Port, reg *mmio.U32) bitband.Bit {
	return bitband.Alias32(reg).Bit(pnum(p))
}
