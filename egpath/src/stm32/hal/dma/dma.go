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

// Channel returns n-th channel (1 <= n <= number of channels).
func (p *DMA) Channel(n int) *Channel {
	return getChannel(p, n)
}

type Channel struct {
	chanregs
}

type Events byte

const (
	GE  = 1 << 0 // Global Event.
	TCE = 1 << 1 // Transfer Complete Event.
	HCE = 1 << 2 // Half transfer Complete Event.
	ERR = 1 << 3 // Transfer Error event.
)

// Events returns current event flags.
func (c *Channel) Events() Events {
	return events(c)
}

// ClearEvents clears specified event flags.
func (c *Channel) ClearEvents(e Events) {
	clearEvents(c, e)
}

// Enable enables channel.
func (c *Channel) Enable() {
	enable(c)
}

// Disable disables channel.
func (c *Channel) Disable() {
	disable(c)
}

// IntEnabled returns true if at least one from events (excluding GE) can
// generate interrupt.
func (c *Channel) IntEnabled(e Events) bool {
	return intEnabled(c, e)
}

// EnableInt enables interrupt generation by events (excluding GE).
func (c *Channel) EnableInt(e Events) {
	enableInt(c, e)
}

// DisableInt disables interrupt generation by events.
func (c *Channel) DisableInt(e Events) {
	disableInt(c, e)
}

type Mode uint16

const (
	ReadP Mode = 0 << 4 // Read from Peripheral, write to memory.
	ReadM Mode = 1 << 4 // Read from Memory, write to peripheral.
	Circ  Mode = 1 << 5 // Enable circular mode.
	IncP  Mode = 1 << 6 // Peripheral increment mode.
	IncM  Mode = 1 << 7 // Memory increment mode.

	PrioL Mode = 0 << 12 // Channel priority level: Low.
	PrioM Mode = 1 << 12 // Channel priority level: Medium.
	PrioH Mode = 2 << 12 // Channel priority level: High.
	PrioV Mode = 3 << 12 // Channel priority level: Very high.

	MTM = 1 << 14 //  Memory to memory mode.
)

// Mode returns current mode of operation.
func (c *Channel) Mode() Mode {
	return mode(c)
}

// SetMode sets mode of operation.
func (c *Channel) SetMode(m Mode) {
	setMode(c, m)
}

// WordSize returns the current word size (in bytes) for peripheral and memory
// side of transfer.
func (c *Channel) WordSize() (p, m uintptr) {
	return wordSize(c)
}

// SetWordSize sets the word size (in bytes) for peripheral and memory side of
// transfer.
func (c *Channel) SetWordSize(p, m uintptr) {
	setWordSize(c, p, m)
}

// Num returns current number of words to transfer.
func (c *Channel) Num() int {
	return num(c)
}

// SetNum sets number of words to transfer (n <= 65535).
func (c *Channel) SetNum(n int) {
	setNum(c, n)
}

// SetAddrP sets peripheral address.
func (c *Channel) SetAddrP(a unsafe.Pointer) {
	setAddrP(c, a)
}

// SetAddrM sets memory address.
func (c *Channel) SetAddrM(a unsafe.Pointer) {
	setAddrM(c, a)
}
