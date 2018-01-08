package main

import (
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
	dig     [8]gpio.Pins // 0-3 top display, 4-7 bottom display.
	seg     [8]gpio.Pins // A B C D E F G :
	digAll  gpio.Pins
	segAll  gpio.Pins
	dl      [8]gpio.Pins
	n, last byte
}

func (d *Display) SetDigPin(digit int, pin gpio.Pins) {
	d.dig[digit] = pin
}

func (d *Display) SetSegPin(segment int, pin gpio.Pins) {
	d.seg[segment] = pin
}

func (d *Display) SetupPins() {
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

func (d *Display) ISR() {
	pins := d.dl[d.n]
	if d.n < d.last {
		d.n++
	} else {
		d.n = 0
	}
	p0.SetPins(d.digAll)
	p0.ClearPins(d.segAll)
	p0.ClearPins(pins & d.digAll)
	p0.SetPins(pins & d.segAll)
}

//emgo:const
var digits = [16]byte{
	// 0, 1
}

func (d *Display) PrintChar(pos, c byte) {

}

func (d *Display) PrintInt(pos, i int) {

}
