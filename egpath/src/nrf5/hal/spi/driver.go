package spi

import (
	"rtos"
)

// Driver is interrupt based driver to the SPI peripheral.
type Driver struct {
	P *Periph

	txptr uintptr
	txmax uintptr
	rxptr uintptr
	rxmax uintptr
	n     int
	isr   func(*Driver)
	done  rtos.EventFlag
	w     uint16
	swp   int8 // Used for 16-bit transfer: 1 - swap bytes, 0 - do not swap.
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph) *Driver {
	d := new(Driver)
	d.P = p
	return d
}

func (d *Driver) Enable() {
	p := d.P
	d.swp = int8(p.LoadCONFIG()&LSBF) ^ 1 // nRF5 is little-endian.
	p.StoreENABLE(true)
}

func (d *Driver) Disable() {
	p := d.P
	p.StoreENABLE(false)
}

// ISR should be used as SPI interrupt handler.
func (d *Driver) ISR() {
	d.isr(d)
}

// Wait waits for the end of SPI transaction. It must be called after any Async*
// method to ensure that the started transaction has been finished. If Wait is
// called after any of AsyncWrite*Read* methods it returns the number of
// bytes/words read.
func (d *Driver) Wait() int {
	d.done.Wait(1, 0)
	d.P.Event(READY).DisableIRQ()
	return d.n
}
