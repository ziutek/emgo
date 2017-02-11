package adc

import (
	"rtos"
	"sync/fence"

	"stm32/hal/dma"
)

type DriverError byte

const ErrTimeout DriverError = 1

func (e DriverError) Error() string {
	switch e {
	case ErrTimeout:
		return "timeout"
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
func NewDriver(p *Periph, dma *dma.Channel) *Driver {
	d := new(Driver)
	d.P = p
	d.DMA = dma
	return d
}

func (d *Driver) ISR() {
	if p := d.P; p.Event()&d.watch != 0 {
		p.DisableIRQ(EvAll)
		d.done.Signal(1)
	}
}

func (d *Driver) SetDeadline(deadline int64) {
	d.deadline = deadline
}

// ClearAndWatch combined with Wait can be used to wait for any from events. See
// Wait for more information. It is low level function, intended to help to use
// d.P directly.
func (d *Driver) ClearAndWatch(events Event) {
	d.done.Reset(0)
	d.watch = events
	p := d.P
	p.Clear(events)
	p.EnableIRQ(events)
	fence.W() // To order writes to normal and I/O memory.
}

// Wait waits for any from events setup by ClearAndWatch. It returns true if
// event occured or false in case of timeout. It is low level function, intended
// to help to use d.P directly:
//	d.ClearAndWatch(events)
//	startSomething(d.P)
//	if !d.Wait() {
//		// Timeout
//	}
func (d *Driver) Wait() bool {
	return d.done.Wait(1, d.deadline)
}

// Enable enables ADC and waits for Ready event.
func (d *Driver) Enable() error {
	d.ClearAndWatch(Ready)
	d.P.Enable()
	if !d.Wait() {
		return ErrTimeout
	}
	return nil
}
