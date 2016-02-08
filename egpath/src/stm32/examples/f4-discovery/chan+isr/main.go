// This example shows how to use channels to divide interrupt handler into two
// parts: fast part - that runs in interrupt context and slow part - that runs
// in user context.
//
// Fast part only enqueues events/data that receives from source (you) and
// informs the source (using blue LED) if its buffer is full. Slow part
// dequeues events/data and handles them. This scheme can be used to implement
// interrupt driven I/O library.
package main

import (
	"delay"
	"rtos"

	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var LED *gpio.Port

const (
	Green  = gpio.Pin12
	Orange = gpio.Pin13
	Red    = gpio.Pin14
	Blue   = gpio.Pin15

	Button = gpio.Pin0
)

func init() {
	system.Setup168(8)
	systick.Setup()

	gpio.A.EnableClock(true)
	bport := gpio.A
	gpio.D.EnableClock(false)
	LED = gpio.D

	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	LED.Setup(Green|Orange|Red|Blue, cfg)

	// Setup external interrupt source: user button.
	bport.Setup(Button, &gpio.Config{Mode: gpio.In})
	line := exti.Lines(Button)
	line.Connect(bport)
	line.EnableRiseTrig()
	line.EnableInt()

	rtos.IRQ(irq.EXTI0).Enable()
}

var (
	c   = make(chan gpio.Pins, 3)
	led = Green
)

func buttonISR() {
	exti.Lines(Button).ClearPending()
	select {
	case c <- led:
		if led <<= 1; led > Red {
			led = Green
		}
	default:
		// Signal that c is full.
		LED.SetPins(Blue)
		delay.Loop(1e5)
		LED.ClearPins(Blue)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI0: buttonISR,
}

func toggle(leds gpio.Pins) {
	LED.SetPins(leds)
	delay.Millisec(500)
	LED.ClearPins(leds)
	delay.Millisec(500)
}

func main() {
	for {
		toggle(<-c)
	}
}
