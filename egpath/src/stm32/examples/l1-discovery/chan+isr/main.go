// This example shows how to use channels to divide interrupt handler into two
// parts: fast part - that runs in interrupt context and slow part - that runs
// in thread context.
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

var leds *gpio.Port

const (
	Blue  = gpio.Pin6
	Green = gpio.Pin7

	Button = gpio.Pin0
)

func init() {
	system.Setup32(0)
	systick.Setup()

	gpio.A.EnableClock(true)
	bport := gpio.A
	gpio.B.EnableClock(false)
	leds = gpio.B

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Blue, &cfg)

	// Setup external interrupt source: user button.
	bport.Setup(Button, &gpio.Config{Mode: gpio.In})
	line := exti.Lines(Button)
	line.Connect(bport)
	line.EnableRiseTrig()
	line.EnableIRQ()

	rtos.IRQ(irq.EXTI0).Enable()
}

var c = make(chan struct{}, 3)

func buttonISR() {
	exti.Lines(Button).ClearPending()
	select {
	case c <- struct{}{}:
	default:
		// Signal that c is full.
		leds.SetPins(Blue)
		delay.Loop(1e5)
		leds.ClearPins(Blue)
	}
}

func main() {
	for {
		<-c
		leds.SetPins(Green)
		delay.Millisec(500)
		leds.ClearPins(Green)
		delay.Millisec(500)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI0: buttonISR,
}
