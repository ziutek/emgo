package button

import (
	"nrf5/input"

	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/rtc"
)

// IntDrv uses GPIOTE peripheral and interrupts to track changes in a button
// state. Optionally, it can use RTC to implement digital debouncing filter.
// Changes are reported using channel of events. IntDrv is rather expensive by
// counting the resources used to handla only one button but provides energy
// efficient alternative to PollDrv.
type IntDrv struct {
	ch    chan<- input.Event
	src   uint32
	te    gpiote.Chan
	rtc   *rtc.Periph
	delay uint16
	ccn   byte
	val   byte
}

// NewIntDrv returns new IntDrv. It configures pin as input to te GPIOTE
// channel. PullUp determines whether the internal pull-up resistor is used. If
// rtc is non-nil it uses rtc.CC[ccn] to implement digital debouncing. RTC must
// be started before (usually, a free channel of system timer is used). IntDrv
// sends int events to ch with source set to src.
func NewIntDrv(pin gpio.Pin, te gpiote.Chan, pullUp bool, rtc *rtc.Periph, ccn int, ch chan<- input.Event, src uint32) *IntDrv {
	d := new(IntDrv)
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

func (d *IntDrv) handlePinChange() {
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
	case d.ch <- input.MakeIntEvent(d.src, val):
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

// ISR should be called int GPIOTE interrupt handler.
func (d *IntDrv) ISR() {
	if ev := d.te.IN().Event(); ev.IsSet() {
		ev.DisableIRQ()
		d.handlePinChange()
	}
}

// RTCISR should be called int RTC interrupt handler.
func (d *IntDrv) RTCISR() {
	if ev := d.rtc.Event(rtc.COMPARE(int(d.ccn))); ev.IsSet() {
		ev.DisableIRQ()
		d.handlePinChange()
	}
}
