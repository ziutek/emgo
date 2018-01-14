// This code drives ventilation unit controller.
//
// Ventilation unit consist of: a counter-flow recuperator, two DC fans (EBM
// R1G225-AF33-12) and two air filters.
//
// Controller produces two PWM signals to contoll both fans. The current fan
// speed is mesured by counting pulses from speed output. The desired speed can
// be set using rotary encoder. Controller has two 4-digit 7-segment displays,
// that shows the current speed of both fans.
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
	btn     *button.IntDrv
	inputCh = make(chan input.Event, 4)
	pwm     *ppipwm.Toggle
	tach    *Tachometer
	aux     gpio.Pin
)

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC1, 0)

	// Allocate pins (always do it in one place to avoid conflicts).

	p0 := gpio.P0
	disp.SetSegPin(F, gpio.Pin0) // F
	disp.SetDigPin(5, gpio.Pin1) // Bottom 1
	disp.SetDigPin(6, gpio.Pin2) // Bottom 2
	disp.SetSegPin(D, gpio.Pin3) // D
	encBt := p0.Pin(4)
	encA := p0.Pin(5)
	encB := p0.Pin(7)
	// gpio.Pin9 // Left connector
	// gpio.Pin10 // Left connector
	disp.SetSegPin(E, gpio.Pin11) // E
	pwmR := p0.Pin(12)            // Right PWM output.
	disp.SetDigPin(7, gpio.Pin13) // Bottom 3
	tachR := p0.Pin(14)           // Right tach input.
	aux = p0.Pin(15)              // AUX right
	disp.SetDigPin(2, gpio.Pin17) // Top 2
	disp.SetSegPin(G, gpio.Pin18) // G
	disp.SetSegPin(Q, gpio.Pin21) // :
	disp.SetSegPin(B, gpio.Pin22) // B
	disp.SetSegPin(C, gpio.Pin23) // C
	disp.SetDigPin(0, gpio.Pin24) // Top 0
	disp.SetDigPin(3, gpio.Pin25) // Top 3
	disp.SetSegPin(A, gpio.Pin28) // A
	disp.SetDigPin(1, gpio.Pin29) // Top 1
	disp.SetDigPin(4, gpio.Pin30) // Bottom 0

	// Configure pins.

	disp.Setup()
	enc = encoder.New(encA, encB, true, true, inputCh, Encoder)
	btn = button.NewIntDrv(encBt, gpiote.Chan(0), true, rtc.RTC1, 1, inputCh, Button)

	pwm = ppipwm.NewToggle(timer.TIMER1)
	pwm.SetFreq(6, 400) // Gives freq. 1/(400 µs) = 2.5 kHz, PWMmax = 99.
	pwm.Setup(0, pwmR, gpiote.Chan(1), ppi.Chan(0), ppi.Chan(1))

	tach = MakeTachometer(timer.TIMER2, tachR, gpiote.Chan(2), ppi.Chan(2), ppi.Chan(3))

	aux.Setup(gpio.ModeIn)

	// Configure interrupts.

	rtos.IRQ(irq.QDEC).Enable()
	rtos.IRQ(irq.GPIOTE).Enable()
}

func main() {
	n := 0
	rpm := 0
	max := pwm.Max() * 2
	pwm.SetInvVal(0, 0)
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
			pwm.SetInvVal(0, n/2)
			disp.WriteDec(4, 7, 2, n/2)
		default:
			rpm = (rpm*15 + tach.RPM()) / 16
			disp.WriteDec(0, 3, 4, rpm)
			disp.Refresh()
			delay.Millisec(2)
		}
	}
}

func qdecISR() {
	enc.ISR()
}

func gpioteISR() {
	btn.ISR()
}

func rtcISR() {
	rtcst.ISR()
	btn.RTCISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC1:   rtcISR,
	irq.QDEC:   qdecISR,
	irq.GPIOTE: gpioteISR,
}
