package gpiote

import (
	"nrf5/hal/gpio"
	"nrf5/hal/te"

	"nrf5/hal/internal/psel"
)

// Chan represents GPIOTE channel. There are 8 channels numbered from 0 to 7.
type Chan byte

// OUT returns task for writing to pin associated with channel c.
func (c Chan) OUT() *te.Task {
	return r().Regs.Task(int(c))
}

// SET returns task for set pin associated with channel c.
func (c Chan) SET() *te.Task {
	return r().Regs.Task(int(c) + 12)
}

// CLR returns task for clear pin associated with channel c.
func (c Chan) CLR() *te.Task {
	return r().Regs.Task(int(c) + 24)
}

// IN returns event generated from pin associated with channel c.
func (c Chan) IN() *te.Event {
	return r().Regs.Event(int(c))
}

// PORT returns event generated from multiple input pins with SENSE mechanism
// enabled. nRF52.
func PORT() *te.Event {
	return r().Regs.Event(31)
}

type Mode byte

const (
	Disabled Mode = 0 // Disconnect pin from GPIOTE module.
	Event    Mode = 1 // Pin generates IN event.
	Task     Mode = 3 // Pin controlled by OUT, SET, CLR task.
)

type Polarity byte

const (
	None   Polarity = 0 // No task on pin, no event from pin.
	LoToHi Polarity = 1 // OUT task sets pin, IN event when rising edge.
	HiToLo Polarity = 2 // OUT task clears pin, IN event when falling edge.
	Toggle Polarity = 3 // OUT task toggles pin, IN event when any change.
)

type Config struct {
	Mode     Mode
	Polarity Polarity
	OutInit  byte
}

// Config returns current configuration of channel c.
func (c Chan) Config() (gpio.Pin, Config) {
	v := r().config[c].Load()
	return psel.Pin(v >> 8 & 0x3f), Config{
		Mode:     Mode(v & 3),
		Polarity: Polarity(v >> 16 & 3),
		OutInit:  byte(v >> 20 & 1),
	}

}

// Setup setups channel c to use pin and cfg configuration.
func (c Chan) Setup(pin gpio.Pin, cfg Config) {
	v := psel.Sel(pin)<<8 | uint32(cfg.Mode) |
		uint32(cfg.Polarity)<<16 | uint32(cfg.OutInit)&1<<20
	r().config[c].Store(v)
}
