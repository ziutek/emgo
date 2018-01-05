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
	"rtos"

	"nrf5/encoder"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var (
	disp Display
	enc  *encoder.Driver
)

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// Allocate pins (always do it in one place to avoid conflicts).

	p0 := gpio.P0
	disp.SetSegPin(F, gpio.Pin0) // F
	disp.SetDigPin(5, gpio.Pin1) // Bottom 1
	disp.SetDigPin(6, gpio.Pin2) // Bottom 2
	disp.SetSegPin(D, gpio.Pin3) // D
	encBtn := p0.Pin(4)
	encA := p0.Pin(5)
	encB := p0.Pin(7)
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
	enc = encoder.New(encA, encB, encBtn, gpiote.Chan(0), true)

	// Configure interrupts.
	rtos.IRQ(irq.QDEC).Enable()
	rtos.IRQ(irq.GPIOTE).Enable()
}

func main() {
	p0 := gpio.P0
	p0.ClearPins(disp.dig[0])
	p0.ClearPins(disp.segAll)
	p0.SetPins(disp.seg[A])
	n := 0
	for ev := range enc.Events() {
		n += ev.Offset()
		m := n / 2 % 6
		if m < 0 {
			m = 6 + m
		}
		pins := disp.seg[m]
		if ev.Button() {
			pins |= disp.seg[G]
		}
		p0.ClearPins(disp.segAll)
		p0.SetPins(pins)
	}
}

func qdecISR() {
	enc.QDECISR()
}

func gpioteISR() {
	enc.GPIOTEISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:   rtcst.ISR,
	irq.QDEC:   qdecISR,
	irq.GPIOTE: gpioteISR,
}
