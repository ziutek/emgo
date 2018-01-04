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

type Encoder struct {
	a  gpio.Pins
	b  gpio.Pins
	bt gpio.Pins
}

var (
	disp Display
	enc  Encoder
)

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// Allocate pins (always do it in one place to avoid conflicts).

	disp.SetSegPin(F, gpio.Pin0) // F
	disp.SetDigPin(5, gpio.Pin1) // Bottom 1
	disp.SetDigPin(6, gpio.Pin2) // Bottom 2
	disp.SetSegPin(D, gpio.Pin3) // D
	enc.bt = gpio.Pin4
	enc.a = gpio.Pin5
	enc.b = gpio.Pin7
	disp.SetSegPin(E, gpio.Pin9)  // E
	disp.SetDigPin(7, gpio.Pin11) // Bottom 3
	disp.SetDigPin(2, gpio.Pin15) // Top 2
	disp.SetSegPin(G, gpio.Pin17) // G
	disp.SetSegPin(B, gpio.Pin22) // B
	disp.SetSegPin(C, gpio.Pin23) // C
	disp.SetSegPin(Q, gpio.Pin21) // :
	disp.SetDigPin(0, gpio.Pin24) // Top 0
	disp.SetDigPin(3, gpio.Pin25) // Top 3
	disp.SetSegPin(A, gpio.Pin28) // A
	disp.SetDigPin(1, gpio.Pin29) // Top 1
	disp.SetDigPin(4, gpio.Pin30) // Bottom 0

	// Configure pins.

	disp.SetupPins()
	gpio.P0.Setup(enc.a|enc.b|enc.bt, gpio.ModeIn|gpio.PullUp)
}

func main() {
	p0 := gpio.P0
	p0.ClearPins(disp.dig[0])
	for {
		in := p0.Load()
		p0.ClearPins(disp.segAll)
		var out gpio.Pins
		if in&enc.a != 0 {
			out |= disp.seg[F]
		}
		if in&enc.b != 0 {
			out |= disp.seg[B]
		}
		if in&enc.bt != 0 {
			out |= disp.seg[D]
		}
		p0.SetPins(out)
		delay.Millisec(10)
	}

	p0.SetPins(disp.seg[A] | disp.seg[B] | disp.seg[C] | disp.seg[D] | disp.seg[G])
	for {
		for _, dig := range disp.dig {
			p0.SetPins(disp.digAll)
			p0.ClearPins(dig)
			delay.Millisec(2)
		}
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
