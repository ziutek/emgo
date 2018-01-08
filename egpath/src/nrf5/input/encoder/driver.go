// Package encoder provides simple, interrupt driven, driver to rotary encoder.
package encoder

import (
	"nrf5/input"

	"nrf5/hal/gpio"
	"nrf5/hal/qdec"
)

// Driver uses QDEC peripheral to track changes in an encoder position. Changes
// are reported using channel of events.
type Driver struct {
	ch  chan<- input.Event
	src byte
}

// New returns new Driver. It configures a, b as input pins for QDEC peripheral.
// PullUp determines whether the internal pull-up resistors are connected.
// Driver sends events to ch with source set to src.
func New(a, b gpio.Pin, pull, dbf bool, ch chan<- input.Event, src byte) *Driver {
	cfg := gpio.ModeIn
	if pull {
		cfg |= gpio.PullUp
	}
	a.Setup(cfg)
	b.Setup(cfg)
	qd := qdec.QDEC
	if dbf {
		qd.StoreDBFEN(true)
	}
	qd.StorePSEL(qdec.A, a)
	qd.StorePSEL(qdec.B, b)
	qd.StoreSAMPLEPER(qdec.P1ms)
	qd.StoreREPORTPER(qdec.P40)
	qd.StoreSHORTS(qdec.REPORTRDY_READCLRACC)
	qd.StoreENABLE(true)
	qd.Task(qdec.START).Trigger()
	qd.Event(qdec.REPORTRDY).EnableIRQ()
	d := new(Driver)
	d.ch = ch
	d.src = src
	return d
}

func (d *Driver) ISR() {
	qd := qdec.QDEC
	qd.Event(qdec.REPORTRDY).Clear()
	select {
	case d.ch <- input.MakeEvent(d.src, qd.LoadACCREAD()):
	default:
	}
}
