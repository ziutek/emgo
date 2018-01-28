package main

import (
	"sync/atomic"

	"nrf5/hal/gpio"
	"nrf5/hal/rtc"
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
	rtc    *rtc.Periph
	delay  uint16
	ccn    byte
	n      byte
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

// Refresh display next, not empty symbol from internal display list. It returns
// index of dieplayed element.
func (d *Display) Refresh() int {
	var pins gpio.Pins
	n := int(d.n)
	for {
		pins = gpio.Pins(atomic.LoadUint32((*uint32)(&d.dl[n])))
		if n++; n == len(d.dl) {
			n = 0
		}
		if pins != 0 || n == int(d.n) {
			break
		}
	}
	p0 := gpio.P0
	p0.SetPins(d.digAll)
	p0.ClearPins(d.segAll)
	p0.ClearPins(pins & d.digAll)
	p0.SetPins(pins & d.segAll)
	d.n = byte(n)
	return n
}

// Clear clears symbol in display list at address addr. Empty element of display
// list is not displayed so it does not consume time during refresh.
func (d *Display) Clear(addr int) {
	atomic.StoreUint32((*uint32)(&d.dl[addr]), 0)
}

// WriteSym wites symbol to the display list at address addr. Symbol consists of
// segments specified by segs and will be displayed at position pos.
func (d *Display) WriteSym(addr, pos int, segs byte) {
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
	1<<F | 1<<E,                             // I
	1<<E | 1<<D | 1<<C | 1<<B,               // J
	0,
	1<<F | 1<<E | 1<<D, // L
	0,
	1<<E | 1<<G | 1<<C,               // n
	1<<G | 1<<E | 1<<D | 1<<C,        // o
	1<<F | 1<<E | 1<<A | 1<<B | 1<<G, // P
	1<<F | 1<<A | 1<<B | 1<<G | 1<<C, // q
	1<<E | 1<<G,                      // r
	1<<A | 1<<F | 1<<G | 1<<C | 1<<D, // S
	1<<F | 1<<E | 1<<D | 1<<G,        // t
	1<<E | 1<<D | 1<<C,               // u
	0,
	0,
	0,
	1<<F | 1<<G | 1<<B | 1<<C | 1<<D, // y
}

// WriteChar writes symbol to the display list that corresponds to the ASCII
// character c. Symbol will be written at address addr and displayed at position
// pos. Unsupported characters are written as space.
func (d *Display) WriteChar(addr, pos int, c byte) {
	switch {
	case '0' <= c && c <= '9':
		c = digits[c-'0']
	case 'a' <= c && c <= 'y':
		c = letters[c-'a']
	case 'A' <= c && c <= 'Y':
		c = letters[c-'A']
	case c == '-':
		c = 1 << G
	case c == '_':
		c = 1 << D
	case c == '=':
		c = 1<<D | 1<<G
	default:
		c = 0
	}
	d.WriteSym(addr, pos, c)
}

// WriteDigit symbol to the display list that corresponds to hexadecimal digit
// digit. Symbol will be written at address addr and displayed at position pos.
// Digit <0 or >15 will be displayed as spece.
func (d *Display) WriteDigit(addr, pos, digit int) {
	var c byte
	switch {
	case 0 <= digit && digit <= 9:
		c = digits[digit]
	case 10 <= digit && digit <= 16:
		c = letters[digit-10]
	}
	d.WriteSym(addr, pos, c)
}

// WriteNumber writes width symbols to the display list that corresponds to
// signed number. Symbols will be written starting from address addr. The least
// significant digit will be displayed at position pos.
func (d *Display) WriteNumber(addr, pos, width, number, base int) {
	neg := number < 0
	if neg {
		width--
		number = -number
	}
	i := 0
	for ; i < width; i++ {
		if number == 0 {
			if i == 0 {
				d.WriteDigit(addr, pos, 0)
				i++
			}
			break
		}
		d.WriteDigit(addr+i, pos-i, number%base)
		number /= base
	}
	if number != 0 {
		for i = 0; i < width; i++ {
			d.WriteSym(addr+i, pos-i, 1<<A|1<<D)
		}
	}
	if neg {
		d.WriteSym(addr+i, pos-i, 1<<G)
		i++
		width++
	}
	for i < width {
		d.WriteSym(addr+i, pos-i, 0)
		i++
	}
}

// WriteDec is shorthand for WriteNumber(addr, pos, width, number, 10).
func (d *Display) WriteDec(addr, pos, width, number int) {
	d.WriteNumber(addr, pos, width, number, 10)
}

// WriteHex is shorthand for WriteNumber(addr, pos, width, number, 16).
func (d *Display) WriteHex(addr, pos, width, number int) {
	d.WriteNumber(addr, pos, width, number, 16)
}

func (d *Display) WriteString(addr, pos, width int, s string) {
	k := 0
	for i := 0; i < width; i++ {
		if k < len(s) {
			d.WriteChar(addr+i, pos+i, s[k])
			k++
		} else {
			d.WriteSym(addr+i, pos+i, 0)
		}
	}
}

// UseRTC setups rtc.CC[ccn] compare register to generate interrupts with a
// period periodms millisecond. RTC should be started before (usually, a free
// channel of system timer is used). RTC interrupt handling must be enabled in
// NVIC to do not miss a compare event.
func (d *Display) UseRTC(rt *rtc.Periph, ccn, periodms int) {
	d.rtc = rt
	d.ccn = byte(ccn)
	d.delay = uint16(32768 * uint32(periodms) / ((rt.LoadPRESCALER() + 1) * 1e3))
	ev := rt.Event(rtc.COMPARE(int(d.ccn)))
	ev.Clear()
	rt.StoreCC(ccn, rt.LoadCOUNTER()+uint32(d.delay))
	ev.EnableIRQ()
}

// RTCISR should be called int RTC interrupt handler. It checks the compare
// event flag and if set it calls Refresh and updates compare register to
// generate next event.
func (d *Display) RTCISR() int {
	rt, ccn := d.rtc, int(d.ccn)
	if ev := rt.Event(rtc.COMPARE(int(ccn))); ev.IsSet() {
		ev.Clear()
		n := d.Refresh()
		rt.StoreCC(ccn, rt.LoadCOUNTER()+uint32(d.delay))
		return n
	}
	return -1
}
