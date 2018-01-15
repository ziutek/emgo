// This code drives ventilation unit controller.
//
// Ventilation unit consist of: a counter-flow recuperator, two DC fans (EBM
// R1G225-AF33-12) and two air filters.
//
// Controller produces two PWM signals to contoll both fans. The current fan
// speed is mesured by counting pulses from TACH output. The desired speed can
// be set using rotary encoder. Controller has two 4-digit 7-segment displays,
// that shows the current speed of both fans.
//
// Work in progress...
package main

import (
	"delay"
	"rtos"

	"nrf5/input"
	"nrf5/input/button"
	"nrf5/input/encoder"
	"nrf5/ppipwm"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/irq"
	"nrf5/hal/ppi"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/timer"
)

const (
	Encoder = iota
	Button
)

var (
	disp    Display
	enc     *encoder.Driver
	btn     *button.PollDrv
	inputCh = make(chan input.Event, 4)
	pwm     *ppipwm.Toggle
	tach    [2]*Tachometer
	aux     gpio.Pin
)

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC1, 0)

	// Allocate pins (always do it in one place to avoid conflicts).

	p0 := gpio.P0
	disp.SetSegPin(F, gpio.Pin0)  // Segment F.
	disp.SetDigPin(5, gpio.Pin1)  // Bottom digit 1.
	disp.SetDigPin(6, gpio.Pin2)  // Bottom digit 2.
	disp.SetSegPin(D, gpio.Pin3)  // Segment D.
	encBtn := gpio.Pin4           // Encoder push button.
	encA := p0.Pin(5)             // Encoder A-phase input.
	encB := p0.Pin(7)             // Encoder B-phase input.
	tach0 := p0.Pin(9)            // Left tach input.
	pwm0 := p0.Pin(10)            // Left PWM output.
	disp.SetSegPin(E, gpio.Pin11) // Segment E.
	pwm1 := p0.Pin(12)            // Right PWM output.
	disp.SetDigPin(7, gpio.Pin13) // Bottom digit 3.
	tach1 := p0.Pin(14)           // Right tach input.
	aux = p0.Pin(15)              // Right AUX.
	disp.SetDigPin(2, gpio.Pin17) // Top digit 2.
	disp.SetSegPin(G, gpio.Pin18) // Segment G.
	disp.SetSegPin(Q, gpio.Pin21) // Segment :.
	disp.SetSegPin(B, gpio.Pin22) // Segment B.
	disp.SetSegPin(C, gpio.Pin23) // Segment C.
	disp.SetDigPin(0, gpio.Pin24) // Top digit 0.
	disp.SetDigPin(3, gpio.Pin25) // Top digit 3.
	disp.SetSegPin(A, gpio.Pin28) // Segment A.
	disp.SetDigPin(1, gpio.Pin29) // Top digit 1.
	disp.SetDigPin(4, gpio.Pin30) // Bottom 0.

	// Configure pins.

	disp.Setup()
	disp.UseRTC(rtc.RTC1, 1, 3)

	enc = encoder.New(encA, encB, true, true, inputCh, Encoder)

	btn = button.NewPollDrv(p0, encBtn, true, inputCh, Button)
	btn.UseRTC(rtc.RTC1, 2, 20)

	pwm = ppipwm.NewToggle(timer.TIMER1)
	pwm.SetFreq(6, 400) // Gives freq. 1/(400 µs) = 2.5 kHz, PWMmax = 99.
	pwm.Setup(0, pwm0, gpiote.Chan(0), ppi.Chan(0), ppi.Chan(1))
	pwm.Setup(1, pwm1, gpiote.Chan(1), ppi.Chan(2), ppi.Chan(3))

	tach[0] = MakeTachometer(
		timer.TIMER2, tach0, gpiote.Chan(2), ppi.Chan(4), ppi.Chan(5),
	)
	tach[1] = MakeTachometer(
		timer.TIMER3, tach1, gpiote.Chan(3), ppi.Chan(6), ppi.Chan(7),
	)

	aux.Setup(gpio.ModeIn)

	// Configure interrupts.

	rtos.IRQ(irq.QDEC).Enable()
}

func main() {
	n := 0
	rpm := 0
	max := pwm.Max() * 2
	pwm.SetInvVal(1, 0)
	disp.WriteDec(4, 7, 2, 0)
	for i := 0; ; i++ {
		select {
		case ev := <-inputCh:
			switch ev.Src() {
			case Encoder:
				n += ev.Int()
				switch {
				case n < 0:
					n = 0
				case n > max:
					n = max
				}
			case Button:
				n = 0
			}
			pwm.SetInvVal(1, n/2)
			disp.WriteDec(4, 7, 2, n/2)
		default:
			rpm = (rpm*15 + tach[1].RPM()) / 16
			disp.WriteDec(0, 3, 4, rpm)
			delay.Millisec(40)
		}
	}
}

func qdecISR() {
	enc.ISR()
}

func rtcISR() {
	rtcst.ISR()
	disp.RTCISR()
	btn.RTCISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC1: rtcISR,
	irq.QDEC: qdecISR,
}
