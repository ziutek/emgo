// Package encoder provides simple, interrupt driven, driver to rotary encoder
// with button.
package encoder

import (
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/qdec"
)

// Event represents change in encoder state.
type Event struct {
	val int
}

// Offset returns the offset from previous position of a rotor.
func (ev Event) Offset() int {
	return ev.val >> 1
}

// Button returns current state of an encoder button.
func (ev Event) Button() bool {
	return ev.val&1 == 0
}

// Driver uses QDEC peripheral and one GPIOTE channel to track changes in
// the rotary encoder state. Changes are reported using channel of events.
type Driver struct {
	ch  chan Event
	btn gpiote.Chan
}

// New returns new Driver. It configures a, b, btn as input pins. PullUp
// determines whether the internal pull-up resistors are connected. It also
// configures QDEC peripheral and uses a, b pins as source of phase-A and
// phase-B signals. Additionally it configures GPIOTE te channel to detect
// changes at btn pin. Set btn to gpio.Pin{} if button is not used.
func New(a, b, btn gpio.Pin, te gpiote.Chan, pullUp bool) *Driver {
	d := new(Driver)
	d.ch = make(chan Event, 3)
	cfg := gpio.ModeIn
	if pullUp {
		cfg |= gpio.PullUp
	}
	a.Setup(cfg)
	b.Setup(cfg)
	if btn.IsValid() {
		btn.Setup(cfg)
		te.Setup(btn, gpiote.ModeEvent|gpiote.PolarityToggle)
		te.IN().Event().EnableIRQ()
		d.btn = te
	} else {
		d.btn = -1
	}
	qd := qdec.QDEC
	qd.StorePSEL(qdec.A, a)
	qd.StorePSEL(qdec.B, b)
	qd.StoreSAMPLEPER(qdec.P1ms)
	qd.StoreREPORTPER(qdec.P40)
	qd.StoreSHORTS(qdec.REPORTRDY_READCLRACC)
	qd.StoreDBFEN(true)
	qd.StoreENABLE(true)
	qd.Task(qdec.START).Trigger()
	qd.Event(qdec.REPORTRDY).EnableIRQ()

	return d
}

func (d *Driver) Events() <-chan Event {
	return d.ch
}

func (d *Driver) QDECISR() {
	qd := qdec.QDEC
	qd.Event(qdec.REPORTRDY).Clear()
	ev := Event{qd.LoadACCREAD() << 1}
	if d.btn >= 0 {
		pin, _ := d.btn.Config()
		ev.val |= pin.Load()
	}
	select {
	case d.ch <- ev:
	default:
	}
}

func (d *Driver) GPIOTEISR() {
	d.btn.IN().Event().Clear()
	pin, _ := d.btn.Config()
	select {
	case d.ch <- Event{pin.Load()}:
	default:
	}
}
