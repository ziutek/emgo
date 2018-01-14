package button

import (
	"nrf5/input"

	"nrf5/hal/gpio"
	"nrf5/hal/rtc"
)

// PollDrv is polling driver that can handle multiple buttons connected to the
// same GPIO port. Optionally, it can use one RTC compare register to generate
// periodic interrupts that are used to call its Poll method. The new button
// state is reported only if it is stable during two subsequent polls.
type PollDrv struct {
	ch    chan<- input.Event
	src   uint32
	port  *gpio.Port
	pins  gpio.Pins
	val   gpio.Pins
	prev  gpio.Pins
	rtc   *rtc.Periph
	delay uint16
	ccn   byte
}

// NewPollDrv returns new PollDrv. It can handle multiple pins of the given
// port. It configures all pins as inputs with pull-up enabled if pullUp is
// true. PollDrv sends pin events to ch with source set to src.
func NewPollDrv(port *gpio.Port, pins gpio.Pins, pullUp bool, ch chan<- input.Event, src uint32) *PollDrv {
	d := new(PollDrv)
	d.ch = ch
	d.src = src
	d.port = port
	d.pins = pins
	cfg := gpio.ModeIn
	if pullUp {
		cfg |= gpio.PullUp
	}
	port.Setup(pins, cfg)
	d.val = d.port.Pins(d.pins)
	d.prev = d.val
	return d
}

// Poll should be caled periodically to poll the state of pins.
func (d *PollDrv) Poll() {
	cur := d.port.Pins(d.pins)
	// Only bits (pins) with the same value in cur and d.prev are trated as
	// stable and will replace bits in d.val.
	mask := cur ^ ^d.prev
	d.prev = cur
	val := d.val&^mask | cur&mask
	if val == d.val {
		return
	}
	d.val = val
	select {
	case d.ch <- input.MakePinEvent(d.src, val):
	default:
	}
}

// UseRTC setups rtc.CC[ccn] compare register to generate interrupts with a
// period periodms millisecond. RTC should be started before (usually, a free
// channel of system timer is used). RTC interrupt handling must be enabled in
// NVIC to do not miss a compare event.
func (d *PollDrv) UseRTC(rt *rtc.Periph, ccn, periodms int) {
	d.rtc = rt
	d.ccn = byte(ccn)
	d.delay = uint16(32768 * uint32(periodms) / ((rt.LoadPRESCALER() + 1) * 1e3))
	ev := rt.Event(rtc.COMPARE(int(d.ccn)))
	ev.Clear()
	rt.StoreCC(ccn, rt.LoadCOUNTER()+uint32(d.delay))
	ev.EnableIRQ()
}

// RTCISR should be called int RTC interrupt handler. It checks the compare
// event flag and if set it calls Poll and updates compare register to generate
// next event.
func (d *PollDrv) RTCISR() {
	rt, ccn := d.rtc, int(d.ccn)
	if ev := rt.Event(rtc.COMPARE(int(ccn))); ev.IsSet() {
		ev.Clear()
		d.Poll()
		rt.StoreCC(ccn, rt.LoadCOUNTER()+uint32(d.delay))
	}
}
