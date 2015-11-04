// Package dcf77 decodes DCF77 time signal.
package dcf77

import (
	"fmt"
	"time"
)

type Error int

const (
	ErrInit   = Error(-1)
	ErrTiming = Error(-2)
	ErrBits   = Error(-3)
)

var strerr = []string{
	"initializing",
	"timing error",
	"bits error",
}

func (e Error) Error() string {
	i := uint(-e - 1)
	if i >= uint(len(strerr)) {
		return "unknown"
	}
	return strerr[i]
}

type Date struct {
	Year   int8
	Month  int8
	Mday   int8
	Wday   int8
	Hour   int8
	Min    int8
	Sec    int8
	Summer bool
}

func (t Date) Format(f fmt.State, _ rune) {
	zone := "CET"
	if t.Summer {
		zone = "CES"
	}
	fmt.Fprintf(
		f,
		"%02d-%02d-%02d %02d:%02d:%02d %s",
		t.Year, t.Month, t.Mday, t.Hour, t.Min, t.Sec, zone,
	)
}

type pulse struct {
	stamp time.Time
	l     uint32
	h     uint16
	sec   int8 // If sec < 0 thne sec can be only ErrInit or ErrTiming.
}

type Decoder struct {
	// ISR fields.
	pulse pulse
	n     byte

	// User fields.
	date Date

	// Common fields.
	c chan pulse
}

// NewDecoder returns pointer to new ready to use DCF77 signal decoder.
func NewDecoder() *Decoder {
	d := new(Decoder)
	d.pulse.sec = int8(ErrInit)
	d.c = make(chan pulse, 1)
	return d
}

func checkRising(dt64 time.Duration) int {
	if dt64 > 2050e6 {
		return -1
	}
	dt := uint(dt64)
	switch {
	case dt > 1950e6:
		return 1
	case dt > 1050e6:
		return -1
	case dt > 950e6:
		return 0
	}
	return -1
}

func (d *Decoder) risingEdge(dt time.Duration) {
	switch checkRising(dt) {
	case 0: // Ordinary pulse.
		d.n++
		if d.pulse.sec >= 0 {
			d.pulse.sec = int8(d.n)
		}
	case 1: // Sync pulse.
		d.n = 0
		if d.pulse.sec >= int8(ErrInit) {
			d.pulse.sec = 0
		} else {
			d.pulse.sec = int8(ErrInit)
		}
	default:
		d.pulse.sec = int8(ErrTiming)
	}
}

func checkFalling(dt64 time.Duration) int {
	if dt64 > 250e6 {
		return -1
	}
	dt := uint(dt64)
	switch {
	case dt > 140e6:
		return 1
	case dt > 130e6:
		return -1
	case dt > 40e6:
		return 0
	}
	return -1
}

func (d *Decoder) fallingEdge(dt time.Duration) {
	if Error(d.pulse.sec) == ErrTiming {
		return
	}
	bit := checkFalling(dt)
	if bit < 0 {
		d.pulse.sec = int8(ErrTiming)
	}
	n := int(d.n) - 16
	switch {
	case n < 0:
		// Skip weather and antena info.
	case n == 0:
		d.pulse.l = uint32(bit)
		d.pulse.h = 0
	case n < 32:
		d.pulse.l += uint32(bit << uint(n))
	default:
		d.pulse.h += uint16(bit << uint(n-32))
	}
}

// Edge should be called by interrupt handler trigered by both (rising and
// falling) edges of DCF77 signal pulses.
func (d *Decoder) Edge(t time.Time, rising bool) {
	dt := t.Sub(d.pulse.stamp)
	lastsec := d.pulse.sec
	if rising {
		d.pulse.stamp = t
		d.risingEdge(dt)
	} else {
		d.fallingEdge(dt)
	}
	if d.pulse.sec != lastsec {
		select {
		case d.c <- d.pulse:
		default:
		}
	}
}

type Pulse struct {
	Date
	Stamp time.Time
}

func (p *Pulse) Err() error {
	if e := Error(p.Date.Sec); e < 0 {
		return e
	}
	return nil
}

func checkParity(u, pbit uint32) bool {
	return true
}

func decodeBCD(u uint32) (int8, bool) {
	h := u >> 4 & 0x0f
	l := u & 0x0f
	return int8(h*10 + l), l < 10 // Don't check h because result is always checked.
}

func (d *Decoder) decodeDate(l, h uint32) {
	ok := true
	switch l >> (17 - 16) & 3 {
	case 2:
		d.date.Summer = false
	case 1:
		d.date.Summer = true
	default:
		ok = false
	}
	ok = ok && l&(1<<(20-16)) != 0
	var o bool
	u := l >> (21 - 16) & 0x7f
	d.date.Min, o = decodeBCD(u)
	ok = ok && o && d.date.Min < 60 && checkParity(u, l>>(28-16))
	u = l >> (29 - 16) & 0x3f
	d.date.Hour, o = decodeBCD(u)
	ok = ok && o && d.date.Hour < 24 && checkParity(u, l>>(35-16))

	u = l>>(36-16) + h<<(32-36+16)
	d.date.Mday, o = decodeBCD(u >> (36 - 36) & 0x3f)
	ok = ok && o && uint(d.date.Mday)-1 < 31
	d.date.Wday = int8(l >> (42 - 16))
	ok = ok && d.date.Wday != 0
	d.date.Month, o = decodeBCD(u >> (45 - 36) & 0x1f)
	ok = ok && o && uint(d.date.Month)-1 < 12
	d.date.Year, o = decodeBCD(u >> (50 - 36))
	ok = ok && o && uint(d.date.Year) < 100
	ok = ok && checkParity(u&0x3fffff, u>>22)
	if ok {
		d.date.Sec = 0
	} else {
		d.date.Sec = int8(ErrBits)
	}
}

// Pulse returns next decoded pulse. Decoder contains internal buffer for one
// value, so if Pulse is called with period > 1 second, it should be called
// twice to obtain most recent value.
func (d *Decoder) Pulse() Pulse {
	var p pulse
	for {
		p = <-d.c
		if p.sec == 0 {
			d.decodeDate(p.l, uint32(p.h))
			break
		}
		if d.date.Sec >= 0 || p.sec < 0 {
			d.date.Sec = p.sec
			break
		}
	}
	return Pulse{d.date, p.stamp}
}
