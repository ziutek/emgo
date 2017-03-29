package gpiote

import (
	"mmio"
	"unsafe"

	"nrf5/hal/gpio"
	"nrf5/hal/te"

	"nrf5/hal/internal/mmap"
	"nrf5/hal/internal/psel"
)

type Periph struct {
	te.Regs

	_      [68]mmio.U32
	config [8]mmio.U32
}

//emgo:const
var GPIOTE = (*Periph)(unsafe.Pointer(mmap.BaseAPB + 0x06000))

const badTask = "gpio: bad task"

// OUT returns task for writing to pin associated with channel n.
func (p *Periph) OUT(n int) *te.Task {
	if uint(n) > 7 {
		panic(badTask)
	}
	return p.Regs.Task(n)
}

// SET returns task for set pin associated with channel n.
func (p *Periph) SET(n int) *te.Task {
	if uint(n) > 7 {
		panic(badTask)
	}
	return p.Regs.Task(n + 12)
}

// CLR returns task for clear pin associated with channel n.
func (p *Periph) CLR(n int) *te.Task {
	if uint(n) > 7 {
		panic(badTask)
	}
	return p.Regs.Task(n + 24)
}

// IN returns event generated from pin associated with channel n.
func (p *Periph) IN(n int) *te.Event {
	return p.Regs.Event(n)
}

// PORT returns event generated from multiple input pins with SENSE mechanism
// enabled.
func (p *Periph) PORT() *te.Event {
	return p.Regs.Event(31)
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
	Pin      gpio.Pin
	Mode     Mode
	Polarity Polarity
	OutInit  byte
}

// LoadCONFIG returns current configuration of n-th channel.
func (p *Periph) LoadCONFIG(n int) Config {
	c := p.config[n].Load()
	return Config{
		Pin:      psel.Pin(c >> 8 & 0x3f),
		Mode:     Mode(c & 3),
		Polarity: Polarity(c >> 16 & 3),
		OutInit:  byte(c >> 20 & 1),
	}
}

// StoreCONFIG stores configuration for n-th channel.
func (p *Periph) StoreCONFIG(n int, cfg Config) {
	c := psel.Sel(cfg.Pin)<<8 | uint32(cfg.Mode) |
		uint32(cfg.Polarity)<<16 | uint32(cfg.OutInit)&1<<20
	p.config[n].Store(c)
}
