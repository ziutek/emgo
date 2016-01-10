package gpio

// Port represents GPIO port.
type Port struct {
	registers
}

// Num returns ordinar number of p on its peripheral bus.
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
	enr().ClearBit(pnum(p))
}

// Reset resets port p.
func (p *Port) Reset() {
	n := pnum(p)
	rstr().SetBit(n)
	rstr().ClearBit(n)
}
