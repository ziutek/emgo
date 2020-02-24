package adc

import (
	"bits"
	"rtos"
	"sync/atomic"
	"sync/fence"
	"unsafe"

	"stm32/hal/dma"
)

type DriverError byte

const (
	ErrDrvOverrun DriverError = 1
)

func (e DriverError) Error() string {
	switch e {
	case ErrDrvOverrun:
		return "drv.overrun"
	}
	return ""
}

type Driver struct {
	P   *Periph
	DMA *dma.Channel

	done    rtos.EventFlag
	waitfor uint32
	offset  byte
}

// NewDriver provides convenient way to create heap allocated Driver.
func NewDriver(p *Periph, ch *dma.Channel) *Driver {
	d := new(Driver)
	d.P = p
	d.DMA = ch
	return d
}

func (d *Driver) ISR() {
	waitfor := d.waitfor
	if waitfor == 0 {
		// Other ADC (shared IRQ).
		return
	}
	p := d.P
	if ev, err := p.Status(); (uint32(ev)<<16|uint32(err))&waitfor == 0 {
		// Other ADC (shared IRQ).
		return
	}
	p.DisableIRQ(EvAll, ErrAll)
	d.waitfor = 0
	d.done.Signal(1)
}

// Enable enables ADC.
func (d *Driver) Enable(calibrate bool) {
	d.enable(calibrate)
}

func (d *Driver) DMAISR() {
	d.DMA.DisableIRQ(dma.EvAll, dma.ErrAll)
	d.done.Signal(1)
}

func (d *Driver) watch(ev Event, err Error) {
	d.done.Reset(0)
	if waitfor := uint32(ev)<<16 | uint32(err); waitfor != 0 {
		d.waitfor = uint32(ev)<<16 | uint32(err)
		p := d.P
		p.Clear(ev, err)
		p.EnableIRQ(ev, err)
	}
	fence.W() // To order writes to normal and I/O memory.
}

func (d *Driver) readDMA(maddr unsafe.Pointer, n int, wsize uintptr) (int, error) {
	if n > 0xffff {
		n = 0xffff
	}
	p, ch := d.P, d.DMA
	paddr := p.raw.DR.U32.Addr()
	if wsize == 1 {
		paddr += uintptr(d.offset)
	}
	enableDMA(ch, 0, 0, unsafe.Pointer(paddr), maddr, wsize, n)
	p.EnableDMA(false)
	d.watch(0, ErrAll)
	p.Start()
	d.done.Wait(1, 0)
	p.Stop()
	ch.Disable() // Required by STM32F1 to allow setup next transfer.
	p.DisableDMA()
	var err error
	switch {
	case ErrAll != 0 && atomic.LoadUint32(&d.waitfor) == 0:
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
	d.offset = byte(bits.One(rdmsb))
}

// Read uses DMA in one shot mode so can not read more than 65535 samples.
func (d *Driver) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	return d.readDMA(unsafe.Pointer(&buf[0]), len(buf), 1)
}

// Read16 uses DMA in one shot mode so can not read more than 65535 samples.
func (d *Driver) Read16(buf []uint16) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	return d.readDMA(unsafe.Pointer(&buf[0]), len(buf), 2)
}
