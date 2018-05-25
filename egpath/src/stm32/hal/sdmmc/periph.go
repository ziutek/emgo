package sdmmc

import (
	"unsafe"

	"stm32/hal/internal"
	"stm32/hal/system"
)

type Periph periph

// Bus returns a bus to which p is connected to.
func (p *Periph) Bus() system.Bus {
	return internal.Bus(unsafe.Pointer(p))
}

// EnableClock enables clock for p.
// lp determines whether the clock remains on in low power (sleep) mode.
func (p *Periph) EnableClock(lp bool) {
	addr := unsafe.Pointer(p)
	internal.APB_SetLPEnabled(addr, lp)
	internal.APB_SetEnabled(addr, true)
}

// DisableClock disables clock for p.
func (p *Periph) DisableClock() {
	internal.APB_SetEnabled(unsafe.Pointer(p), false)
}

// Reset resets p.
func (p *Periph) Reset() {
	internal.APB_Reset(unsafe.Pointer(p))
}

// Enabled reports whether the p is enabled.
func (p *Periph) Enabled() bool {
	return p.raw.PWRCTRL().Load() == 3
}

// Enable enables p. At least seven PCLK clock periods are needed between any
// Enable or Disable.
func (p *Periph) Enable() {
	p.raw.POWER.Store(3)
}

// Disable disables gp. At least seven PCLK clock periods are needed between any
// Enable or Disable.
func (p *Periph) Disable() {
	p.raw.POWER.Store(0)
}

type BusClock byte

const (
	ClkEna   BusClock = 1 << 0 // Enables bus clock.
	PwrSave  BusClock = 1 << 1 // Enables power saving mode.
	ClkByp   BusClock = 1 << 2 // Pass SDMMC clock directly to CK pin.
	BusWidth BusClock = 3 << 3 // Describes data bus width.
	Bus1     BusClock = 0 << 3 // Single data bus line.
	Bus4     BusClock = 1 << 3 // Four data bus lines.
	Bus8     BusClock = 2 << 3 // Eight data bus lines.
	NegEdge  BusClock = 1 << 5 // Command and data changed on CK falling edge.
	FlowCtrl BusClock = 1 << 6 // Enables hardware flow controll.
)

// BusClock returns the current configuration of SDMMC bus.
func (p *Periph) BusClock() (cfg BusClock, clkdiv uint) {
	clkcr := p.raw.CLKCR.Load()
	return BusClock(clkcr >> 8), uint(clkdiv & 255)
}

// SetBusClock configures the SDMMC bus.
//
// Note (Errata Sheet DocID027036 Rev 2):
// 2.7.1 Don't use HW flow control (FlowCtrl).
// 2.7.3 Don't use clock dephasing (NegEdge).
// 2.7.5 Ensure:
//
//	3*period(PCLK)+3*period(SDMMCCLK) < 32/BusWidth*period(SDMMC_CK)
//  always met for: PCLK > 28.8 MHz).
//
func (p *Periph) SetBusClock(cfg BusClock, clkdiv uint) {
	if clkdiv > 255 {
		panic("sdio: bad clkdiv")
	}
	p.raw.CLKCR.U32.Store(uint32(cfg)<<8 | uint32(clkdiv))
}
