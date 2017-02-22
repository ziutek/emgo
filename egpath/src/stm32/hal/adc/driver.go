package adc

import (
	"bits"
	"rtos"
	"sync/fence"
	"unsafe"

	"stm32/hal/dma"
)

type DriverError byte

const ErrTimeout DriverError = 1

func (e DriverError) Error() string {
	switch e {
	case ErrTimeout:
		return "timeout"
	}
	return ""
}

type Driver struct {
	deadline int64

	P   *Periph
	DMA *dma.Channel

	done       rtos.EventFlag
	watch      uint32
	byteOffset byte
}

// NewDriver provides convenient way to create heap allocated Driver.
func NewDriver(p *Periph, ch *dma.Channel) *Driver {
	d := new(Driver)
	d.P = p
	d.DMA = ch
	return d
}

func (d *Driver) ISR() {
	if d.watch == 0 {
		// Other peripheral (shared IRQ).
		return
	}
	p := d.P
	if ev, err := p.Status(); (uint32(err)<<16|uint32(ev))&d.watch == 0 {
		// Other peripheral (shared IRQ).
		return
	}
	p.DisableIRQ(EvAll, ErrAll)
	d.watch = 0
	d.done.Signal(1)
}

func (d *Driver) SetDeadline(deadline int64) {
	d.deadline = deadline
}

// Watch combined with Wait can be used to wait for any from events/errors. See
// Wait for more information. It is low level function, intended to help to use
// d.P directly.
func (d *Driver) Watch(ev Event, err Error) {
	d.done.Reset(0)
	d.watch = uint32(err)<<16 | uint32(ev) | 1<<15 // Assumes unused 15 bit.
	p := d.P
	p.Clear(ev, err)
	p.EnableIRQ(ev, err)
	fence.W() // To order writes to normal and I/O memory.
}

// Wait waits for any from events/error setup by Watch or for DMA events/errors.
// It returns true if event/error occured or false in case of timeout. It is
// low level function, intended to help to use d.P directly:
//	d.Watch(events, errors)
//	startSomething(d.P)
//	if !d.Wait() {
//		// Timeout
//	}
func (d *Driver) Wait() bool {
	return d.done.Wait(1, d.deadline)
}

// Happened returns true if an event/error set by Watch occured.
func (d *Driver) Happened() bool {
	return d.watch == 0
}

// Enable enables ADC.
func (d *Driver) Enable(calibrate bool) error {
	return d.enable(calibrate)
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
	paddr := p.raw.DR.U32.Addr()
	if wsize == 1 {
		paddr += uintptr(d.byteOffset)
	}
	setupDMA(ch, paddr, wsize)
	startDMA(ch, addr, n)
	p.EnableDMA(false)
	d.Watch(0, ErrAll)
	acceptTrig(p)
	timeout := !d.Wait()
	ch.Disable() // Required by STM32F1 to allow setup next transfer.
	p.DisableDMA()
	var err error
	switch {
	case timeout:
		ch.DisableIRQ(dma.EvAll, dma.ErrAll)
		p.DisableIRQ(EvAll, ErrAll)
		err = ErrTimeout
	case d.Happened():
		ch.DisableIRQ(dma.EvAll, dma.ErrAll)
		_, err = p.Status()
	default:
		p.DisableIRQ(EvAll, ErrAll)
		if _, e := ch.Status(); e&^dma.ErrFIFO != 0 {
			err = e
		}
	}
	return n - ch.Len(), err
}

// SetReadMSB sets most significant byte of 16-bit ADC data register to be read
// by Read and ReadByte methods.
func (d *Driver) SetReadMSB(rdmsb bool) {
	d.byteOffset = byte(bits.One(rdmsb))
}

// Read uses DMA in one shot mode so can not read more than 65535 samples.
func (d *Driver) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	return d.readDMA(uintptr(unsafe.Pointer(&buf[0])), len(buf), 1)
}

// Read16 uses DMA in one shot mode so can not read more than 65535 samples.
func (d *Driver) Read16(buf []uint16) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	return d.readDMA(uintptr(unsafe.Pointer(&buf[0])), len(buf), 2)
}
