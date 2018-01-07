// Package button provides interrupt driven driver to push button.
package button

import (
	"nrf5/input"

	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/rtc"
)

// Driver uses GPIOTE peripheral to track changes in a button position.
// Optionally, it can use RTC to implement digital debouncing filter. Changes
// are reported using channel of events. Driver is rather expensive by counting
// the resources used but provides energy efficient alternative to polling
// implementation.
type Driver struct {
	ch    chan<- input.Event
	te    gpiote.Chan
	rtc   *rtc.Periph
	delay uint16
	ccn   byte
	src   byte
	val   byte
}

// New returns new Driver. It configures pin as input to te GPIOTE channel.
// PullUp determines whether the internal pull-up resistor is connected. If rtc
// is non-nil it uses rtc.CC[ccn] to implement digital debouncing. RTC must be
// started before (usually, a free channel of system timer is used). Driver
// sends events to ch with source set to src.
func New(pin gpio.Pin, te gpiote.Chan, pullUp bool, rtc *rtc.Periph, ccn int, ch chan<- input.Event, src byte) *Driver {
	d := new(Driver)
	d.ch = ch
	d.src = src
	d.te = te
	d.val = 2
	if rtc != nil {
		d.rtc = rtc
		d.ccn = byte(ccn)
		d.delay = uint16(32768 / 16 / (rtc.LoadPRESCALER() + 1)) // 1/16 s.
	}
	cfg := gpio.ModeIn
	if pullUp {
		cfg |= gpio.PullUp
	}
	pin.Setup(cfg)
	te.Setup(pin, gpiote.ModeEvent|gpiote.PolarityToggle)
	te.IN().Event().EnableIRQ()
	return d
}

func (d *Driver) handlePinChange() {
	te := d.te
	ev := te.IN().Event()
	ev.Clear() // Clear before load the pin value, to don't miss any change.
	pin, _ := te.Config()
	val := pin.Load()
	if val == int(d.val) {
		ev.EnableIRQ()
		return
	}
	d.val = byte(val)
	select {
	case d.ch <- input.MakeEvent(d.src, val):
	default:
	}
	rt := d.rtc
	if rt == nil {
		return
	}
	ccn := int(d.ccn)
	rt.StoreCC(ccn, rt.LoadCOUNTER()+uint32(d.delay))
	ev = rt.Event(rtc.COMPARE(ccn))
	ev.Clear()
	ev.EnableIRQ()
}

func (d *Driver) ISR() {
	d.te.IN().Event().DisableIRQ()
	d.handlePinChange()
}

func (d *Driver) RTCISR() {
	d.rtc.Event(rtc.COMPARE(int(d.ccn))).DisableIRQ()
	d.handlePinChange()
}
