package gpiote

import (
	"nrf5/hal/gpio"
	"nrf5/hal/te"
)

// Chan represents GPIOTE channel. There are 4 (8 in nRF52) channels numbered
// from 0 to 3 (7 in nRF52).
type Chan byte

// OUT returns task for writing to pin associated with channel c.
func (c Chan) OUT() *te.Task {
	return r().Regs.Task(int(c))
}

// SET returns task for set pin associated with channel c. nRF52.
func (c Chan) SET() *te.Task {
	return r().Regs.Task(int(c) + 12)
}

// CLR returns task for clear pin associated with channel c. nRF52.
func (c Chan) CLR() *te.Task {
	return r().Regs.Task(int(c) + 24)
}

// IN returns event generated from pin associated with channel c.
func (c Chan) IN() *te.Event {
	return r().Regs.Event(int(c))
}

// PORT returns event generated from multiple input pins with SENSE mechanism
// enabled.
func PORT() *te.Event {
	return r().Regs.Event(31)
}

type Config uint32

const (
	ModeDiscon Config = 0 // Disconnect pin from GPIOTE module.
	ModeEvent  Config = 1 // Pin generates IN event.
	ModeTask   Config = 3 // Pin controlled by OUT, SET, CLR task.

	PolarityNone   Config = 0 << 16 // No task on pin, no event from pin.
	PolarityLoToHi Config = 1 << 16 // OUT sets pin, IN when rising edge.
	PolarityHiToLo Config = 2 << 16 // OUT clears pin, IN when falling edge.
	PolarityToggle Config = 3 << 16 // OUT toggles pin, IN when any change.

	OutInitLow  Config = 0       // Low initial output value.
	OutInitHigh Config = 1 << 20 // High initial output value.
)

// Config returns current configuration of channel c.
func (c Chan) Config() (gpio.Pin, Config) {
	v := r().config[c].Load()
	const psel = 0x7 << 8
	return gpio.SelPin(int8(v & psel >> 8)), Config(v &^ psel)
}

// Setup setups channel c to use pin and cfg configuration.
func (c Chan) Setup(pin gpio.Pin, cfg Config) {
	r().config[c].Store(uint32(pin.Sel())<<8 | uint32(cfg))
}
