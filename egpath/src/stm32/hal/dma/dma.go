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

// Stream returns n-th stream (channel in F1 nomenclature).
func (p *DMA) Stream(n int) *Stream {
	return getStream(p, n)
}

type Stream struct {
	stregs
}

type Events byte

const (
	TCE = tce // Transfer Complete Event.
	HCE = hce // Half transfer Complete Event.
	ERR = err // Error event.
)

// Events returns current event flags.
func (s *Stream) Events() Events {
	return events(s)
}

// ClearEvents clears specified event flags.
func (s *Stream) ClearEvents(e Events) {
	clearEvents(s, e)
}

// Enable enables channel.
func (s *Stream) Enable() {
	enable(s)
}

// Disable disables channel.
func (s *Stream) Disable() {
	disable(s)
}

// IntEnabled returns events that are enabled to generate interrupts.
func (s *Stream) IntEnabled() Events {
	return intEnabled(s)
}

// EnableInt enables interrupt generation by events.
func (s *Stream) EnableInt(e Events) {
	enableInt(s, e)
}

// DisableInt disables interrupt generation by events.
func (s *Stream) DisableInt(e Events) {
	disableInt(s, e)
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
)

type Channel byte

// Setup configures stream.
func (s *Stream) Setup(m Mode, ch Channel) {
	setup(s, m, ch)
}

// WordSize returns the current word size (in bytes) for peripheral and memory
// side of transfer.
func (s *Stream) WordSize() (p, m uintptr) {
	return wordSize(s)
}

// SetWordSize sets the word size (in bytes) for peripheral and memory side of
// transfer.
func (s *Stream) SetWordSize(p, m uintptr) {
	setWordSize(s, p, m)
}

// Num returns current number of words to transfer.
func (s *Stream) Num() int {
	return num(s)
}

// SetNum sets number of words to transfer (n <= 65535).
func (s *Stream) SetNum(n int) {
	setNum(s, n)
}

// SetAddrP sets peripheral address.
func (s *Stream) SetAddrP(a unsafe.Pointer) {
	setAddrP(s, a)
}

// SetAddrM sets memory address.
func (s *Stream) SetAddrM(a unsafe.Pointer) {
	setAddrM(s, a)
}
