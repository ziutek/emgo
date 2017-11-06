package dma

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	DMA1 = (*DMA)(unsafe.Pointer(mmap.DMA1_BASE))
	DMA2 = (*DMA)(unsafe.Pointer(mmap.DMA2_BASE))
)

type DMA dmaperiph

func (p *DMA) EnableClock(lp bool) {
	p.enableClock(lp)
}

func (p *DMA) DisableClock() {
	p.disableClock()
}

func (p *DMA) Reset() {
	p.reset()
}

// Channel returns value that represents sn-th stream (channel in F1/L1 series
// nomenclature) with cn-th request channel set (ignored in case of F1/L1
// series). Channels with the same sn points to the same DMA stream so they can
// not be used concurently.
func (p *DMA) Channel(sn, cn int) *Channel {
	return p.getChannel(sn, cn)
}

type Channel channel

type Event byte

const (
	Complete     Event = trce // Transfer Complete Event.
	HalfComplete Event = htce // Half Transfer Complete Event.

	EvAll = Complete | HalfComplete
)

type Error byte

const (
	ErrTransfer   Error = trerr // Transfer Error.
	ErrDirectMode Error = dmerr // Direct Mode Error.
	ErrFIFO       Error = fferr // FIFO Error.

	ErrAll = ErrTransfer | ErrDirectMode | ErrFIFO
)

func (e Error) Error() string {
	var (
		s string
		d Error
	)
	switch {
	case e&ErrTransfer != 0:
		d = ErrTransfer
		s = "DMA transfer+"
	case e&ErrDirectMode != 0:
		d = ErrDirectMode
		s = "DMA direct mode+"
	case e&ErrFIFO != 0:
		d = ErrFIFO
		s = "DMA FIFO+"
	default:
		return ""
	}
	if e&^d == 0 {
		s = s[:len(s)-1]
	}
	return s
}

// Status returns current event and error flags.
func (ch *Channel) Status() (Event, Error) {
	flags := ch.status()
	return Event(flags) & EvAll, Error(flags) & ErrAll
}

// ClearEvents clears specified event flags.
func (ch *Channel) Clear(ev Event, err Error) {
	ch.clear(byte(ev) | byte(err))
}

// Enable enables the channel ch. All events and errors should be cleared
// before call this method.
func (ch *Channel) Enable() {
	ch.enable()
}

// Disable disables channel.
func (ch *Channel) Disable() {
	ch.disable()
}

// Returns true if channel is enabled.
func (ch *Channel) Enabled() bool {
	return ch.enabled()
}

// IRQEnabled returns events that are enabled to generate interrupt requests.
func (ch *Channel) IRQEnabled() (Event, Error) {
	flags := ch.irqEnabled()
	return Event(flags) & EvAll, Error(flags) & ErrAll
}

// EnableIRQ enables generation of IRQs by ev, err. Documentation does not
// mention it, but IRQ can be not generated if an event was asserted before
// enable IRQ for it. So always enable IRQs before channel. Typically, the
// correct sequence is as follows:
//	ch.Clear(EvAll, ErrAll)
//	ch.EnableIRQ(ev, err)
//	ch.Enable()
func (ch *Channel) EnableIRQ(ev Event, err Error) {
	ch.enableIRQ(byte(ev) | byte(err))
}

// DisableIRQ disables IRQs generation by ev, err.
func (ch *Channel) DisableIRQ(ev Event, err Error) {
	ch.disableIRQ(byte(ev) | byte(err))
}

type Mode uint32

const (
	PTM Mode = 0   // Read from peripheral, write to memory.
	MTP Mode = mtp // Read from memory, write to peripheral.
	MTM Mode = mtm // Read from memory (AddrP), write to memory.

	Circ Mode = circ // Enable circular mode.
	IncP Mode = incP // Peripheral increment mode.
	IncM Mode = incM // Memory increment mode.

	Direct   Mode = 0        // Direct mode.
	FIFO_1_4 Mode = fifo_1_4 // FIFO mode, threshold 1/4.
	FIFO_2_4 Mode = fifo_2_4 // FIFO mode, threshold 2/4.
	FIFO_3_4 Mode = fifo_3_4 // FIFO mode, threshold 3/4.
	FIFO_4_4 Mode = fifo_4_4 // FIFO mode, threshold 4/4.
)

// Setup configures channel.
func (ch *Channel) Setup(m Mode) {
	ch.setup(m)
}

type Prio byte

const (
	Low      Prio = 0     // Stream priority level: Low.
	Medium   Prio = prioM // Stream priority level: Medium.
	High     Prio = prioH // Stream priority level: High.
	VeryHigh Prio = prioV // Stream priority level: Very high.
)

func (ch *Channel) SetPrio(prio Prio) {
	ch.setPrio(prio)
}

func (ch *Channel) Prio() Prio {
	return ch.prio()
}

// WordSize returns the current word size (in bytes) for peripheral and memory
// side of transfer.
func (ch *Channel) WordSize() (p, m uintptr) {
	return ch.wordSize()
}

// SetWordSize sets the word size (in bytes) for peripheral and memory side of
// transfer.
func (ch *Channel) SetWordSize(p, m uintptr) {
	ch.setWordSize(p, m)
}

// Len returns current number of words to transfer.
func (ch *Channel) Len() int {
	return ch.len()
}

// SetLen sets number of words to transfer (n <= 65535).
func (ch *Channel) SetLen(n int) {
	ch.setLen(n)
}

// SetAddrP sets peripheral address (or memory source address in case of MTM).
func (ch *Channel) SetAddrP(a unsafe.Pointer) {
	ch.setAddrP(a)
}

// SetAddrM sets memory address.
func (ch *Channel) SetAddrM(a unsafe.Pointer) {
	ch.setAddrM(a)
}

// Select selects request source (peripheral) for channel.
func (ch *Channel) Select(req int) {
	ch.sel(req)
}
