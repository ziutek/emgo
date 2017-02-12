package adc

import (
	"rtos"
	"sync/fence"
	"unsafe"

	"stm32/hal/dma"
)

type DriverError byte

const (
	ErrTimeout DriverError = iota
	ErrOverrun
)

func (e DriverError) Error() string {
	switch e {
	case ErrTimeout:
		return "timeout"
	case ErrOverrun:
		return "overrun"
	default:
		return ""
	}
}

type Driver struct {
	deadline int64

	P   *Periph
	DMA *dma.Channel

	done  rtos.EventFlag
	watch Event
}

// NewDriver provides convenient way to create heap allocated Driver.
func NewDriver(p *Periph, ch *dma.Channel) *Driver {
	d := new(Driver)
	d.P = p
	d.DMA = ch
	return d
}

func (d *Driver) ISR() {
	if p := d.P; p.Event()&d.watch != 0 {
		p.DisableIRQ(EvAll)
		d.watch = 0
		d.done.Signal(1)
	}
}

func (d *Driver) SetDeadline(deadline int64) {
	d.deadline = deadline
}

// WatchEvent combined with Wait can be used to wait for any from events. See
// Wait for more information. It is low level function, intended to help to use
// d.P directly.
func (d *Driver) WatchEvent(events Event) {
	d.done.Reset(0)
	d.watch = events
	p := d.P
	p.Clear(events)
	p.EnableIRQ(events)
	fence.W() // To order writes to normal and I/O memory.
}

// WaitEvent waits for any from events setup by WatchEvent or DMA event. It
// returns true if event occured or false in case of timeout. It is low level
// function, intended to help to use d.P directly:
//	d.WatchEvent(events)
//	startSomething(d.P)
//	if !d.WaitEvent() {
//		// Timeout
//	}
func (d *Driver) WaitEvent() bool {
	return d.done.Wait(1, d.deadline)
}

func (d *Driver) EventHappened() bool {
	return d.watch == 0
}

// Enable enables ADC and waits for Ready event.
func (d *Driver) Enable() error {
	d.WatchEvent(Ready)
	d.P.Enable()
	if !d.WaitEvent() {
		return ErrTimeout
	}
	return nil
}

func (d *Driver) DMAISR() {
	d.DMA.DisableIRQ(dma.EvAll, dma.ErrAll)
	d.done.Signal(1)
}

func setupDMA(ch *dma.Channel, paddr uintptr, wordSize uintptr) {
	ch.Setup(dma.PTM | dma.IncM | dma.FIFO_1_4)
	ch.SetWordSize(wordSize, wordSize)
	ch.SetAddrP(unsafe.Pointer(paddr))
}

func startDMA(ch *dma.Channel, addr uintptr, n int) {
	ch.SetAddrM(unsafe.Pointer(addr))
	ch.SetLen(n)
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.EnableIRQ(dma.Complete, dma.ErrAll&^dma.ErrFIFO)
	fence.W() // This orders writes to normal and I/O memory.
	ch.Enable()
}

func (d *Driver) readDMA(addr uintptr, n int, wsize uintptr) (int, error) {
	if n > 0xffff {
		n = 0xffff
	}
	p, ch := d.P, d.DMA
	setupDMA(ch, p.raw.DR.U32.Addr(), wsize)
	startDMA(ch, addr, n)
	p.EnableDMA(false)
	d.WatchEvent(Overrun)
	p.Start()
	timeout := !d.WaitEvent()
	ch.Disable() // Required by STM32F1 to allow setup next transfer.
	p.DisableDMA()
	var err error
	switch {
	case timeout:
		ch.DisableIRQ(dma.EvAll, dma.ErrAll)
		p.DisableIRQ(Overrun)
		err = ErrTimeout
	case d.EventHappened():
		ch.DisableIRQ(dma.EvAll, dma.ErrAll)
		err = ErrOverrun
	default:
		p.DisableIRQ(Overrun)
		if _, e := ch.Status(); e&^dma.ErrFIFO != 0 {
			err = e
		}
	}
	return n - ch.Len(), err
}

// Read uses DMA in one shot mode so can not read more than 65535 samples.
func (d *Driver) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	return d.readDMA(uintptr(unsafe.Pointer(&buf[0])), len(buf), 1)
}
