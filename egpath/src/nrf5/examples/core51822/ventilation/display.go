package main

import (
	"sync/atomic"

	"nrf5/hal/gpio"
)

// Segment names
const (
	A = iota
	B
	C
	D
	E
	F
	G
	Q // Colon
)

// Two 4-digit 7-segment displays (BW428G-E4, common cathode).
type Display struct {
	dig    [8]gpio.Pins // 0-3 top display, 4-7 bottom display.
	seg    [8]gpio.Pins // A B C D E F G :
	digAll gpio.Pins
	segAll gpio.Pins
	dl     [8]gpio.Pins
	n      int
}

func (d *Display) SetDigPin(digit int, pin gpio.Pins) {
	d.dig[digit] = pin
}

func (d *Display) SetSegPin(segment int, pin gpio.Pins) {
	d.seg[segment] = pin
}

func (d *Display) Setup() {
	d.digAll = 0
	d.segAll = 0
	for i := 0; i < 8; i++ {
		d.digAll |= d.dig[i]
		d.segAll |= d.seg[i]
	}
	p0 := gpio.P0
	// Drive digits with higd drive, open drain (n-channel).
	p0.SetPins(disp.digAll)
	p0.Setup(disp.digAll, gpio.ModeOut|gpio.DriveH0D1)
	// Drive segments with higd drive, open drain (p-channel).
	p0.Setup(disp.segAll, gpio.ModeOut|gpio.DriveD0H1)
}

// Refresh display next, not empty, symbol from internal display list.
func (d *Display) Refresh() {
	var pins gpio.Pins
	n := d.n
	for {
		pins = gpio.Pins(atomic.LoadUint32((*uint32)(&d.dl[n])))
		if n++; n == len(d.dl) {
			n = 0
		}
		if pins != 0 || n == d.n {
			break
		}
	}
	d.n = n
	p0 := gpio.P0
	p0.SetPins(d.digAll)
	p0.ClearPins(d.segAll)
	p0.ClearPins(pins & d.digAll)
	p0.SetPins(pins & d.segAll)
}

// Clear clears symbol in display list at address addr. Empty element
// of display list is not displayed so it does not consume time during
// refresh.
func (d *Display) Clear(addr int) {
	atomic.StoreUint32((*uint32)(&d.dl[addr]), 0)
}

// Store stores symbol in display list at address addr. Symbol consists
// of segments specified by segs and will be displayed at position pos.
func (d *Display) Store(addr, pos int, segs byte) {
	pins := d.dig[pos]
	for n := 0; segs != 0; n++ {
		if segs&1 != 0 {
			pins |= d.seg[n]
		}
		segs >>= 1
	}
	atomic.StoreUint32((*uint32)(&d.dl[addr]), uint32(pins))
}

//emgo:const
var digits = [...]byte{
	1<<A | 1<<B | 1<<C | 1<<D | 1<<E | 1<<F,        // 0
	1<<B | 1<<C,                                    // 1
	1<<A | 1<<B | 1<<G | 1<<E | 1<<D,               // 2
	1<<A | 1<<B | 1<<C | 1<<D | 1<<G,               // 3
	1<<F | 1<<G | 1<<B | 1<<C,                      // 4
	1<<A | 1<<F | 1<<G | 1<<C | 1<<D,               // 5
	1<<A | 1<<F | 1<<E | 1<<D | 1<<C | 1<<G,        // 6
	1<<F | 1<<A | 1<<B | 1<<C,                      // 7
	1<<A | 1<<B | 1<<C | 1<<D | 1<<E | 1<<F | 1<<G, // 8
	1<<A | 1<<B | 1<<C | 1<<D | 1<<F | 1<<G,        // 9
}

//emgo:const
var letters = [...]byte{
	1<<E | 1<<F | 1<<A | 1<<B | 1<<C | 1<<G, // A
	1<<F | 1<<E | 1<<D | 1<<C | 1<<G,        // b
	1<<G | 1<<E | 1<<D,                      // c
	1<<B | 1<<E | 1<<D | 1<<C | 1<<G,        // d
	1<<A | 1<<F | 1<<E | 1<<D | 1<<G,        // E
	1<<A | 1<<F | 1<<E | 1<<G,               // F
	1<<A | 1<<F | 1<<E | 1<<D | 1<<C,        // G
	1<<F | 1<<E | 1<<G | 1<<C,               // h
	1 << E, // i
	1<<E | 1<<D | 1<<C | 1<<B, // J
	0,
	1<<F | 1<<E | 1<<D, // L
	0,
	1<<E | 1<<G | 1<<C,               // n
	1<<G | 1<<E | 1<<D | 1<<C,        // o
	1<<F | 1<<E | 1<<A | 1<<B | 1<<G, // P
	1<<F | 1<<A | 1<<B | 1<<G | 1<<C, // q
	1<<E | 1<<G,                      // r
	1<<A | 1<<F | 1<<G | 1<<C | 1<<D, // S
	1<<F | 1<<E | 1<<G,               // t
	1<<E | 1<<D | 1<<C,               // u
	0,
	0,
	0,
	1<<F | 1<<G | 1<<B | 1<<C | 1<<D, // y
}

// StoreChar stores character in display list at address addr. See store
// for more information.
func (d *Display) StoreChar(addr, pos int, c byte) {
	switch {
	case '0' <= c && c <= '9':
		c = digits[c-'0']
	case 'a' <= c && c <= 'y':
		c = letters[c-'a']
	case 'A' <= c && c <= 'Y':
		c = letters[c-'A']
	case c == '-':
		c = 1 << G
	default:
		c = 0
	}
	d.Store(addr, pos, c)
}

// StoreDigit stores digit in display list at address addr. See store
// for more information.
func (d *Display) StoreDigit(addr, pos, digit int) {
	var c byte
	switch {
	case 0 <= digit && digit <= 9:
		c = digits[digit]
	case 10 <= digit && digit <= 19:
		c = letters[digit-10]
	}
	d.Store(addr, pos, c)
}
