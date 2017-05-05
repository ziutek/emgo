package gpiote

import (
	"nrf5/hal/gpio"
	"nrf5/hal/te"
)

// Chan represents GPIOTE channel. There are 4 (8 in nRF52) channels numbered
// from 0 to 3 (7 in nRF52).
type Chan byte

type Task byte

// OUT returns task for writing to pin associated with channel c.
func (c Chan) OUT() Task {
	return Task(c)
}

// SET returns task for set pin associated with channel c. nRF52.
func (c Chan) SET() Task {
	return Task(c + 12)
}

// CLR returns task for clear pin associated with channel c. nRF52.
func (c Chan) CLR() Task {
	return Task(c + 24)
}

type Event byte

const PORT Event = 31 // From multiple input pins with SENSE mechanism enabled.

// IN returns event generated from pin associated with channel c.
func (c Chan) IN() Event {
	return Event(c)
}

func (t Task) Task() *te.Task    { return r().Regs.Task(int(t)) }
func (e Event) Event() *te.Event { return r().Regs.Event(int(e)) }

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
	const psel = 0x7F << 8
	return gpio.SelPin(int8(v & psel >> 8)), Config(v &^ psel)
}

// Setup setups channel c to use pin and cfg configuration.
func (c Chan) Setup(pin gpio.Pin, cfg Config) {
	r().config[c].Store(uint32(pin.Sel())<<8 | uint32(cfg))
}
