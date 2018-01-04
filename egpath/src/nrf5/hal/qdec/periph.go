// Package qdec provides access to the registers of Quadrature Decoder (QDEC)
// peripheral.
package qdec

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

	_          [64]mmio.U32
	enable     mmio.U32
	ledpol     mmio.U32
	sampleper  mmio.U32
	sample     mmio.U32
	reportper  mmio.U32
	acc        mmio.U32
	accread    mmio.U32
	psel       [3]mmio.U32
	dbfen      mmio.U32
	_          [5]mmio.U32
	ledpre     mmio.U32
	accdbl     mmio.U32
	accdblread mmio.U32
}

//emgo:const
var QDEC = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x012000))

type Task byte

const (
	START      Task = 0 // Start the quadrature decoder.
	STOP       Task = 1 // Stop the quadrature decoder.
	READCLRACC Task = 2 // Read and clear ACC and ACCDBL.
	RDCLRACC   Task = 3 // Read and clear ACC. nRF52.
	RDCLRDBL   Task = 4 // Read and clear ACCDBL. nRF52.
)

type Event byte

const (
	SAMPLERDY Event = 0 // New sample value written to the SAMPLE register.
	REPORTRDY Event = 1 // Non-null report ready.
	ACCOF     Event = 2 // ACC or ACCDBL register overflow.
	DBLRDY    Event = 3 // Double displacement(s) detected. nRF52.
	STOPPED   Event = 4 // QDEC has been stopped. nRF52.
)

func (p *Periph) Task(t Task) *te.Task    { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

type Shorts uint32

const (
	REPORTRDY_READCLRACC Shorts = 1 << 0
	SAMPLERDY_STOP       Shorts = 1 << 1
	REPORTRDY_RDCLRACC   Shorts = 1 << 2 // nRF52.
	REPORTRDY_STOP       Shorts = 1 << 3 // nRF52.
	DBLRDY_RDCLRDBL      Shorts = 1 << 4 // nRF52.
	DBLRDY_STOP          Shorts = 1 << 5 // nRF52.
	SAMPLERDY_READCLRACC Shorts = 1 << 6 // nRF52.
)

func (p *Periph) LoadSHORTS() Shorts   { return Shorts(p.Regs.LoadSHORTS()) }
func (p *Periph) StoreSHORTS(s Shorts) { p.Regs.StoreSHORTS(uint32(s)) }

// LoadENABLE reports whether the QDEC peripheral is enabled.
func (p *Periph) LoadENABLE() bool {
	return p.enable.Load()&1 != 0
}

// StoreENABLE enables or disables QDEC peripheral.
func (p *Periph) StoreENABLE(en bool) {
	p.enable.Store(uint32(bits.One(en)))
}

type Signal byte

const (
	LED Signal = 0 // LED output.
	A   Signal = 1 // Phase A input.
	B   Signal = 2 // Phase B input.
)

func (p *Periph) LoadPSEL(s Signal) gpio.Pin {
	return gpio.SelPin(int8(p.psel[s].Load()))
}
func (p *Periph) StorePSEL(s Signal, pin gpio.Pin) {
	p.psel[s].Store(uint32(pin.Sel()))
}
