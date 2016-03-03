package dma

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

var (
	DMA1 = (*DMA)(unsafe.Pointer(mmap.DMA1_BASE))
	DMA2 = (*DMA)(unsafe.Pointer(mmap.DMA2_BASE))
)

type DMA struct {
	dmaregs
}

func (p *DMA) EnableClock(lp bool) {
	enableClock(p, lp)
}

func (p *DMA) DisableClock() {
	disableClock(p)
}

func (p *DMA) Reset() {
	reset(p)
}

// Channel returns value that represents sn-th stream (channel in F1/L1 series
// nomenclature) and cn-th channel (ignored in case of F1/L1 series).
func (p *DMA) Channel(sn, cn int) Channel {
	return getChannel(p, sn, cn)
}

type Channel struct {
	channel
}

type Events byte

const (
	TRCE Events = trce // Transfer Complete Event.
	HTCE Events = htce // Half Transfer Complete Event.
	EV          = TRCE | HTCE

	TRERR Events = trerr // Transfer Error.
	DMERR Events = dmerr // Direct Mode Error.
	FFERR Events = fferr // FIFO Error.
	ERR          = TRERR | DMERR | FFERR
)

// Events returns current event flags.
func (ch Channel) Events() Events {
	return events(ch)
}

// ClearEvents clears specified event flags.
func (ch Channel) ClearEvents(e Events) {
	clearEvents(ch, e)
}

// Enable enables channel.
func (ch Channel) Enable() {
	enable(ch)
}

// Disable disables channel.
func (ch Channel) Disable() {
	disable(ch)
}

// IntEnabled returns events that are enabled to generate interrupts.
func (ch Channel) IntEnabled() Events {
	return intEnabled(ch)
}

// EnableInt enables interrupt generation by events.
func (ch Channel) EnableInt(e Events) {
	enableInt(ch, e)
}

// DisableInt disables interrupt generation by events.
func (ch Channel) DisableInt(e Events) {
	disableInt(ch, e)
}

type Mode uint32

const (
	PTM Mode = 0   // Read from peripheral, write to memory.
	MTP Mode = mtp // Read from memory, write to peripheral.
	MTM Mode = mtm // Read from memory (AddrP), write to memory.

	Circ Mode = circ // Enable circular mode.
	IncP Mode = incP // Peripheral increment mode.
	IncM Mode = incM // Memory increment mode.

	PrioL Mode = 0     // Stream priority level: Low.
	PrioM Mode = prioM // Stream priority level: Medium.
	PrioH Mode = prioH // Stream priority level: High.
	PrioV Mode = prioV // Stream priority level: Very high.

	Direct   = 0        // Direct mode.
	FIFO_1_4 = fifo_1_4 // FIFO mode, threshold 1/4.
	FIFO_2_4 = fifo_2_4 // FIFO mode, threshold 2/4.
	FIFO_3_4 = fifo_3_4 // FIFO mode, threshold 3/4.
	FIFO_4_4 = fifo_4_4 // FIFO mode, threshold 4/4.
)

// Setup configures channel.
func (ch Channel) Setup(m Mode) {
	setup(ch, m)
}

// WordSize returns the current word size (in bytes) for peripheral and memory
// side of transfer.
func (ch Channel) WordSize() (p, m uintptr) {
	return wordSize(ch)
}

// SetWordSize sets the word size (in bytes) for peripheral and memory side of
// transfer.
func (ch Channel) SetWordSize(p, m uintptr) {
	setWordSize(ch, p, m)
}

// Len returns current number of words to transfer.
func (ch Channel) Len() int {
	return length(ch)
}

// SetLen sets number of words to transfer (n <= 65535).
func (ch Channel) SetLen(n int) {
	setLen(ch, n)
}

// SetAddrP sets peripheral address (or memory source address in case of MTM).
func (ch Channel) SetAddrP(a unsafe.Pointer) {
	setAddrP(ch, a)
}

// SetAddrM sets memory address.
func (ch Channel) SetAddrM(a unsafe.Pointer) {
	setAddrM(ch, a)
}
