package spi

import (
	"rtos"

	"arch/cortexm/scb"
)

// Driver is interrupt based driver to the SPI peripheral.
type Driver struct {
	P *Periph

	txbuf string
	txn   int
	rxbuf []byte
	done  rtos.EventFlag
	isr   func(*Driver)
	w16   byte
	swp   int8 // Used for 16-bit transfer: 1 - swap bytes, 0 - do not swap.
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph) *Driver {
	d := new(Driver)
	d.P = p
	return d
}

func (d *Driver) Enable() {
	if d.isr == nil {
		d.isr = (*Driver).isr8
	}
	p := d.P
	cpuBigEndian := int8(scb.SCB.AIRCR.Load() >> scb.ENDIANNESSn & 1)
	d.swp = int8(p.LoadCONFIG()&LSBF) ^ cpuBigEndian
	p.StoreENABLE(true)
	ev := p.Event(READY)
	ev.Clear()
	ev.EnableIRQ()
}

func (d *Driver) Disable() {
	p := d.P
	p.StoreENABLE(false)
	p.Event(READY).DisableIRQ()
}

func (d *Driver) WordSize() int {
	return 8 << d.w16
}

// SetWordSize sets word size in bits (driver supports only 8 and 16 bit).
func (d *Driver) SetWordSize(size int) {
	switch size {
	case 8:
		d.isr = (*Driver).isr8
	case 16:
		d.isr = (*Driver).isr16
	default:
		panic("spi: bad word size")
	}
	d.w16 = byte(size / 16)
}

// ISR should be used as SPI interrupt handler.
func (d *Driver) ISR() {
	d.isr(d)
}

func (d *Driver) Wait() int {
	d.done.Wait(1, 0)
	return len(d.rxbuf) >> d.w16
}
