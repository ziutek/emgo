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
	RDCLRACC   Task = 3 // Read and clear ACC, nRF52.
	RDCLRDBL   Task = 4 // Read and clear ACCDBL, nRF52.
)

type Event byte

const (
	SAMPLERDY Event = 0 // New sample value written to the SAMPLE register.
	REPORTRDY Event = 1 // Non-null report ready.
	ACCOF     Event = 2 // ACC or ACCDBL register overflow.
	DBLRDY    Event = 3 // Double displacement(s) detected, nRF52.
	STOPPED   Event = 4 // QDEC has been stopped, nRF52.
)

func (p *Periph) Task(t Task) *te.Task    { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

type Shorts uint32

const (
	REPORTRDY_READCLRACC Shorts = 1 << 0
	SAMPLERDY_STOP       Shorts = 1 << 1
	REPORTRDY_RDCLRACC   Shorts = 1 << 2 // nRF52
	REPORTRDY_STOP       Shorts = 1 << 3 // nRF52
	DBLRDY_RDCLRDBL      Shorts = 1 << 4 // nRF52
	DBLRDY_STOP          Shorts = 1 << 5 // nRF52
	SAMPLERDY_READCLRACC Shorts = 1 << 6 // nRF52
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

// LoadLEDPOL returns LED output pin polarity: 0 - active low, 1 - active high.
func (p *Periph) LoadLEDPOL() int {
	return int(p.ledpol.Load())
}

// StoreLEDPOL sets LED output pin polarity: 0 - active low, 1 - active high.
func (p *Periph) StoreLEDPOL(polarity int) {
	p.ledpol.Store(uint32(polarity))
}

type SamplePeriod byte

const (
	P128us SamplePeriod = 0  // 128 µs
	P256us SamplePeriod = 1  // 256 µs
	P512us SamplePeriod = 2  // 512 µs
	P1ms   SamplePeriod = 3  // 1024 µs
	P2ms   SamplePeriod = 4  // 2048 µs
	P4ms   SamplePeriod = 5  // 4096 µs
	P8ms   SamplePeriod = 6  // 8192 µs
	P16ms  SamplePeriod = 7  // 16384 µs
	P33ms  SamplePeriod = 8  // 32768 µs, nRF52
	P66ms  SamplePeriod = 9  // 65536 µs, nRF52
	P131ms SamplePeriod = 10 // 131072 µs, nRF52
)

// LoadSAMPLEPER returns the sample period.
func (p *Periph) LoadSAMPLEPER() SamplePeriod {
	return SamplePeriod(p.sampleper.Load())
}

// StoreSAMPLEPER sets the sample period.
func (p *Periph) StoreSAMPLEPER(period SamplePeriod) {
	p.sampleper.Store(uint32(period))
}

// LoadSAMPLE returns the last motion sample: -1, 1 or 2.
func (p *Periph) LoadSAMPLE() int {
	return int(p.sample.Load())
}

type ReportPeriod byte

const (
	P10  ReportPeriod = 0 // 10 samples per report.
	P40  ReportPeriod = 1 // 40 samples per report.
	P80  ReportPeriod = 2 // 80 samples per report.
	P120 ReportPeriod = 3 // 120 samples per report.
	P160 ReportPeriod = 4 // 160 samples per report.
	P200 ReportPeriod = 5 // 200 samples per report.
	P240 ReportPeriod = 6 // 240 samples per report.
	P280 ReportPeriod = 7 // 280 samples per report.
	P1   ReportPeriod = 8 // 1 sample per report. nRF52.
)

// LoadREPORTPER returns the sample period.
func (p *Periph) LoadREPORTPER() ReportPeriod {
	return ReportPeriod(p.reportper.Load())
}

// StoreREPORTPER sets the sample period.
func (p *Periph) StoreREPORTPER(period ReportPeriod) {
	p.reportper.Store(uint32(period))
}

// LoadACC returns the accumulated valid transitions: [-1024..1023].
func (p *Periph) LoadACC() int {
	return int(p.acc.Load())
}

// LoadACCREAD returns the snapshot of the ACC register, updated by the
// READCLRACC or RDCLRACC task.
func (p *Periph) LoadACCREAD() int {
	return int(p.accread.Load())
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

// LoadDBFEN reports whether the input debounce filters are enabled.
func (p *Periph) LoadDBFEN() bool {
	return p.dbfen.Load() != 0
}

// StoreDBFEN enables or disables the input debounce filters.
func (p *Periph) StoreDBFEN(en bool) {
	p.dbfen.Store(uint32(bits.One(en)))
}

// LoadLEDPRE returns the time period the LED is switched on prior to sampling:
// [0..511] µs.
func (p *Periph) LoadLEDPRE() int {
	return int(p.ledpre.Load())
}

// StoreLEDPRE sets the time period the LED is switched on prior to sampling:
// [0..511] µs.
func (p *Periph) StoreLEDPRE(us int) {
	p.ledpre.Store(uint32(us))
}

// LoadACCDBL returns the number of detected double transitions: [0..15].
func (p *Periph) LoadACCDBL() int {
	return int(p.accdbl.Load())
}

// LoadACCDBLREAD returns the snapshot of the ACCDBL register, updated by the
// READCLRACC or RDCLRACC task.
func (p *Periph) LoadACCDBLREAD() int {
	return int(p.accdblread.Load())
}
