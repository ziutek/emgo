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
	"stm32/hal/setup"
)

var LED *gpio.Port

const (
	Green  = 12
	Orange = 13
	Red    = 14
	Blue   = 15
)

const ButtonPin = 0

func init() {
	setup.Performance168(8)

	gpio.A.Enable(false)
	gpio.D.Enable(false)

	LED = gpio.D
	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Orange, gpio.Out)
	LED.SetMode(Red, gpio.Out)
	LED.SetMode(Blue, gpio.Out)

	// Setup external interrupt source: user button.
	bport := gpio.A
	bport.SetMode(ButtonPin, gpio.In)
	line := exti.Line(ButtonPin)
	line.Connect(bport)
	line.EnableRiseTrig()
	line.EnableInt()

	rtos.IRQ(irq.EXTI0).Enable()
}

var (
	c   = make(chan int, 3)
	led = Green
)

func buttonISR() {
	exti.Line(ButtonPin).ClearPending()
	select {
	case c <- led:
		if led++; led > Red {
			led = Green
		}
	default:
		// Signal that c is full.
		LED.SetPin(Blue)
		delay.Loop(1e5)
		LED.ClearPin(Blue)
	}
}

var ISRs = [...]func(){
	irq.EXTI0: buttonISR,
} //c:__attribute__((section(".ISRs")))

func toggle(led int) {
	LED.SetPin(led)
	delay.Millisec(500)
	LED.ClearPin(led)
	delay.Millisec(500)
}

func main() {
	for {
		toggle(<-c)
	}
}
