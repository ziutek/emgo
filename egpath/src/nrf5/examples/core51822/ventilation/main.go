// This code drives ventilation unit controller. Ventilation unit consist of:
// a counter-flow recuperator, two DC fans (EBM R1G225-AF33-12) and two air
// filters.
//
// Controller produces two PWM signals to contoll both fans. The current fan
// speed is mesured by counting pulses from speed output. The desired speed can
// be set using rotary encoder. Controller has two 4-digit 7-segment displays,
// that shows the current speed of both fans.
package main

import (
	"delay"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

const (
	A = iota
	B
	C
	D
	E
	F
	G
	Q
)

// Two 4-digit 7-segment displays (BW428G-E4, common cathode).
type Display struct {
	dig    [8]gpio.Pins // 0-3 top display, 4-7 bottom display.
	seg    [8]gpio.Pins // A B C D E F G :
	digAll gpio.Pins
	segAll gpio.Pins
}

var disp Display

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// Allocate pins (always do it in one place to avoid conflicts).

	p0 := gpio.P0
	disp.seg[F] = gpio.Pin0  // F
	disp.dig[5] = gpio.Pin1  // Bottom 1
	disp.dig[6] = gpio.Pin2  // Bottom 2
	disp.seg[D] = gpio.Pin3  // D
	disp.seg[E] = gpio.Pin9  // E
	disp.dig[7] = gpio.Pin11 // Bottom 3
	disp.dig[2] = gpio.Pin15 // Top 2
	disp.seg[G] = gpio.Pin17 // G
	disp.seg[B] = gpio.Pin22 // B
	disp.seg[C] = gpio.Pin23 // C
	disp.seg[Q] = gpio.Pin21 // :
	disp.dig[0] = gpio.Pin24 // Top 0
	disp.dig[3] = gpio.Pin25 // Top 3
	disp.seg[A] = gpio.Pin28 // A
	disp.dig[1] = gpio.Pin29 // Top 1
	disp.dig[4] = gpio.Pin30 // Bottom 0

	// Configure pins.

	for _, pin := range disp.dig {
		disp.digAll |= pin
	}
	for _, pin := range disp.seg {
		disp.segAll |= pin
	}
	// Drive digits with higd drive, open drain (n-channel).
	p0.SetPins(disp.digAll)
	p0.Setup(disp.digAll, gpio.ModeOut|gpio.DriveH0D1)
	// Drive segments with higd drive, open drain (p-channel).
	p0.Setup(disp.segAll, gpio.ModeOut|gpio.DriveD0H1)
}

func wait() {
	delay.Millisec(500)
}

func main() {
	p0 := gpio.P0
	for {
		for _, dig := range disp.dig {
			p0.SetPins(disp.digAll)
			p0.ClearPins(dig)
			for _, seg := range disp.seg {
				p0.ClearPins(disp.segAll)
				p0.SetPins(seg)
				delay.Millisec(200)
			}
		}
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
