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
	"bufio"
	"bytes"
	//"debug/semihosting"
	//"fmt"
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
	menu    Menu
	enc     *encoder.Driver
	btn     *button.PollDrv
	inputCh = make(chan input.Event, 4)
	fc      *FanControl
	cr0     gpio.Pin
	cl0     gpio.Pin
	cl1     gpio.Pin
	cl2     gpio.Pin
	cb0     gpio.Pin
	cb1     gpio.Pin
)

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC1, 0)

	disp := menu.Display()

	// Allocate pins (always do it in one place to avoid conflicts).

	p0 := gpio.P0
	disp.SetSegPin(F, gpio.Pin0)  // Segment F.
	disp.SetDigPin(5, gpio.Pin1)  // Bottom digit 1.
	disp.SetDigPin(6, gpio.Pin2)  // Bottom digit 2.
	disp.SetSegPin(D, gpio.Pin3)  // Segment D.
	encBtn := gpio.Pin4           // Encoder push button.
	encA := p0.Pin(5)             // Encoder A-phase input.
	cb0 = p0.Pin(6)               // Bottom connector, pin 0.
	encB := p0.Pin(7)             // Encoder B-phase input.
	cb1 = p0.Pin(8)               // Bottom connector, pin 1.
	tach0 := p0.Pin(9)            // Left tach input.
	pwm0 := p0.Pin(10)            // Left PWM output.
	disp.SetSegPin(E, gpio.Pin11) // Segment E.
	pwm1 := p0.Pin(12)            // Right PWM output.
	disp.SetDigPin(7, gpio.Pin13) // Bottom digit 3.
	tach1 := p0.Pin(14)           // Right tach input.
	cl1 = p0.Pin(15)              // Left connector, pin 1.
	cr0 = p0.Pin(16)              // Right connector, pin 0.
	disp.SetDigPin(1, gpio.Pin17) // Top digit 1.
	disp.SetSegPin(G, gpio.Pin18) // Segment G.
	cl0 = p0.Pin(19)              // Left connector pin 0.
	disp.SetDigPin(2, gpio.Pin20) // Top digit 2.
	disp.SetSegPin(Q, gpio.Pin21) // Segment :.
	disp.SetSegPin(B, gpio.Pin22) // Segment B.
	disp.SetSegPin(C, gpio.Pin23) // Segment C.
	disp.SetDigPin(0, gpio.Pin24) // Top digit 0.
	disp.SetDigPin(3, gpio.Pin25) // Top digit 3.
	disp.SetSegPin(A, gpio.Pin28) // Segment A.
	cl2 = p0.Pin(29)              // Left connector, pin 2.
	disp.SetDigPin(4, gpio.Pin30) // Bottom digit 0.

	// Configure pins.

	disp.Setup()
	disp.UseRTC(rtc.RTC1, 1, 3)

	enc = encoder.New(encA, encB, true, true, inputCh, Encoder)

	btn = button.NewPollDrv(p0, encBtn, true, inputCh, Button)
	btn.UseRTC(rtc.RTC1, 2, 20)

	pwm := ppipwm.NewToggle(timer.TIMER1)
	pwm.SetFreq(4, 256) // Gives freq. 1/(256 µs) = 3.9 kHz, PWMmax = 255.
	pwm.Setup(0, pwm0, gpiote.Chan(0), ppi.Chan(0), ppi.Chan(1))
	pwm.SetInv(0, 0) // Immediately stop fan 0.
	pwm.Setup(1, pwm1, gpiote.Chan(1), ppi.Chan(2), ppi.Chan(3))
	pwm.SetInv(1, 0) // Immediately stop fan 1.

	tach := NewTachometer(
		timer.TIMER2, gpiote.Chan(2),
		ppi.Chan(4), ppi.Chan(5), ppi.Chan(6), ppi.Chan(7),
		ppi.Group(0), ppi.Group(1), tach0, tach1,
	)

	fc = NewFanControl(pwm, tach)

	//aux.Setup(gpio.ModeIn)

	// Configure interrupts.

	rtos.IRQ(irq.QDEC).Enable()
	rtos.IRQ(irq.TIMER2).Enable()

	// Semihosting console.

	/*
		f, err := semihosting.OpenFile(":tt", semihosting.W)
		for err != nil {
		}
		fmt.DefaultWriter = lineWriter{bufio.NewWriterSize(f, 40)}
	*/
}

func main() {
	menu.SetMaxRPM(fc.MaxRPM())
	menu.Select(Init)
	fc.Identify()
	menu.Select(ShowRPM)
	for ev := range inputCh {
		switch ev.Src() {
		case Button:
			if ev.Pins() != 0 {
				menu.Next()
			}
		case Encoder:
			n, rpm := menu.HandleEncoder(ev.Int())
			if n >= 0 {
				fc.SetTargetRPM(n, rpm)
			}
		}
	}
}

func qdecISR() {
	enc.ISR()
}

func rtcISR() {
	rtcst.ISR()
	btn.RTCISR()
	menu.RTCISR()
}

func timer2ISR() {
	fc.TachISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC1:   rtcISR,
	irq.QDEC:   qdecISR,
	irq.TIMER2: timer2ISR,
}

type lineWriter struct {
	w *bufio.Writer
}

func (b lineWriter) Write(s []byte) (int, error) {
	n, err := b.w.Write(s)
	if err != nil {
		return n, err
	}
	if bytes.IndexByte(s, '\n') >= 0 {
		err = b.w.Flush()
	}
	return n, err
}
