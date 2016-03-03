package i2c

import (
	"stm32/hal/internal"
	"stm32/hal/system"
	"unsafe"

	"stm32/hal/raw/i2c"
)

type Periph struct {
	raw i2c.I2C_Periph
}

// Bus returns a bus to which p is connected.
func (p *Periph) Bus() system.Bus {
	return internal.Bus(unsafe.Pointer(&p.raw))
}

// EnableClock enables clock for p.
// lp determines whether the clock remains on in low power (sleep) mode.
func (p *Periph) EnableClock(lp bool) {
	addr := unsafe.Pointer(&p.raw)
	internal.APB_SetLPEnabled(addr, lp)
	internal.APB_SetEnabled(addr, true)
}

// DisableClock disables clock for p..
func (p *Periph) DisableClock() {
	internal.APB_SetEnabled(unsafe.Pointer(&p.raw), false)
}

// Reset resets p.
func (p *Periph) Reset() {
	internal.APB_Reset(unsafe.Pointer(&p.raw))
}

func (p *Periph) Enable() {
	p.raw.PE().Set()
}

func (p *Periph) Disable() {
	p.raw.PE().Clear()
}

func (p *Periph) SoftReset() {
	p.raw.SWRST().Set()
	p.raw.SWRST().Clear()
}

// Mode
type Mode byte

const (
	I2C       = Mode(0)
	SMBusDev  = Mode(i2c.SMBUS)
	SMBusHost = Mode(i2c.SMBTYPE | i2c.SMBUS)
)

// DutyCycle describes SCL low time to high time.
type DutyCycle byte

const (
	Duty2_1  DutyCycle = iota // SCL low/high = 2/1.
	Duty16_9                  // SCL low/high = 16/9.
)

type Config struct {
	Speed int       // Clock speed [Hz]
	Duty  DutyCycle // Duty cycle used in fast mode.
	Mode  Mode      // I2C, SMBusDev, SMBusHost.
}

// Setup configures p. It ensures that configured SCL clock speed <= cfg.Speed.
//
// STM32 generates I2C SCL as follow:
//
//   SCL = PCLK / CCR / idiv
//
// where PCLK is peripheral (bus) clock, CCR is 12-bit value from CCR register,
// idiv is internal divider equal to:
//
//   2  for standard mode,
//   3  for fast mode and 2/1 duty cycle,
//   25 for fast mode and 16/9 duty cycle.
//
// For 36 MHz PCLK (typical for 72 MHz F10x) maximum valid SCL can be:
//
//   standard mode:  36 MHz / 180 / 2 = 100 kHz,
//   fast mode 2/1:  36 MHz / 30 / 3  = 400 kHz,
//   fast mode 16/9: 36 MHz / 4 / 25  = 360 kHz.
//
// To obtain 400 kHz SCL in 16/9 fast mode the PCLK must be configured to
// multiple of 10 MHz.
func (p *Periph) Setup(cfg *Config) {
	pclk := int(p.Bus().Clock()) // Pclk should fit in int.
	pclkM := pclk / 1e6
	var ccr, trise int
	if cfg.Speed <= 100e3 {
		// Standard mode.
		div := cfg.Speed * 2
		ccr = (pclk + div - 1) / div
		if ccr < 4 {
			ccr = 4
		}
		trise = pclkM + 1 // SCL max. rise time 1000 ns.
	} else {
		// Fast mode.
		ccr = int(i2c.FS)
		var div int
		if cfg.Duty == Duty2_1 {
			div = cfg.Speed * 3
			ccr = (pclk + div - 1) / div
		} else {
			div = cfg.Speed * 25
			ccr |= int(i2c.DUTY)
		}
		ccrval := (pclk + div - 1) / div
		if ccrval < 1 {
			ccrval = 1
		}
		ccr |= ccrval
		trise = pclkM*3/10 + 1 // SCL max. rise time 300 ns.
	}
	p.raw.CR1.Store(0) // Disables peripheral.
	p.raw.FREQ().Store(i2c.CR2_Bits(pclkM))
	p.raw.CCR.Store(i2c.CCR_Bits(ccr))
	p.raw.TRISE.Store(i2c.TRISE_Bits(trise))
	p.raw.CR1.Store(i2c.CR1_Bits(cfg.Mode))
}

// SPeed returns actual clock speed set.
func (p *Periph) Speed() int {
	ccr := p.raw.CCR.Load()
	var idiv uint
	switch {
	case ccr&i2c.FS == 0: // Standard mode.
		idiv = 2
	case ccr&i2c.DUTY == 0: // Fast mode 2/1.
		idiv = 3
	default: // Fast mode 16/9.
		idiv = 25
	}
	return int(p.Bus().Clock() / (idiv * uint(ccr&i2c.CCRVAL)))
}
