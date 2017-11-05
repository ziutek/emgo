package spi

import (
	"bits"
	"mmio"
	"unsafe"

	"nrf5/hal/gpio"
	"nrf5/hal/te"

	"nrf5/hal/internal/mmap"
)

type Periph struct {
	te.Regs

	_         [64]mmio.U32
	enable    mmio.U32
	_         mmio.U32
	psel      [3]mmio.U32
	_         mmio.U32
	rxd       mmio.U32
	txd       mmio.U32
	_         mmio.U32
	frequency mmio.U32
	_         [11]mmio.U32
	config    mmio.U32
}

//emgo:const
var (
	SPI0 = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x03000))
	SPI1 = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x04000))
)

type Event byte

const READY Event = 2 // TXD byte sent and RXD byte received.

func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

// LoadENABLE reports whether the SPI peripheral is enabled.
func (p *Periph) LoadENABLE() bool {
	return p.enable.Load()&1 != 0
}

// StoreENABLE enables or disables SPI peripheral.
func (p *Periph) StoreENABLE(en bool) {
	p.enable.Store(uint32(bits.One(en)))
}

type Signal byte

const (
	SCK  Signal = 0
	MOSI Signal = 1
	MISO Signal = 2
)

func (p *Periph) LoadPSEL(s Signal) gpio.Pin {
	return gpio.SelPin(int8(p.psel[s].Load()))
}
func (p *Periph) StorePSEL(s Signal, pin gpio.Pin) {
	p.psel[s].Store(uint32(pin.Sel()))
}

func (p *Periph) LoadRXD() byte {
	return byte(p.rxd.Load())
}

func (p *Periph) StoreTXD(b byte) {
	p.txd.Store(uint32(b))
}

// Freq sets
type Freq uint32

const (
	F125k Freq = 0x02000000 // 125 kbps
	F250k Freq = 0x04000000 // 250 kbps
	F500k Freq = 0x08000000 // 500 kbps
	F1M   Freq = 0x10000000 // 1 Mbps
	F2M   Freq = 0x20000000 // 2 Mbps
	F4M   Freq = 0x40000000 // 4 Mbps
	F8M   Freq = 0x80000000 // 8 Mbps
)

// LoadFREQUENCY returns configured SCK frequency.
func (p *Periph) LoadFREQUENCY() Freq {
	return Freq(p.frequency.Load())
}

// StoreFREQUENCY stores SCK frequency.
func (p *Periph) StoreFREQUENCY(f Freq) {
	p.frequency.Store(uint32(f))
}

// Config is a bitfield that describes SPI configuration.
type Config byte

const (
	MSBF  Config = 0      // Most significant bit shifted out first.
	LSBF  Config = 1 << 0 // Most significant bit shifted out first.
	CPHA0 Config = 0      // Sample on leading SCK edge, shift data on trailing edge.
	CPHA1 Config = 1 << 1 // Sample on trailing SCK edge, shift data on leading edge.
	CPOL0 Config = 0      // SCK polarity: active high.
	CPOL1 Config = 1 << 2 // SCK polarity: active low.
)

// LoadCONFIG returns current SPI configuration.
func (p *Periph) LoadCONFIG() Config {
	return Config(p.config.Load())
}

// StoreCONFIG stores SPI configuration..
func (p *Periph) StoreCONFIG(c Config) {
	p.config.Store(uint32(c))
}
