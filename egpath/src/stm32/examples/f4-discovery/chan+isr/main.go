// This example shows how use channels to divide interrupt handler into two
// parts: fast part - that runs in interrupt context and soft part - that run
// in user context. Fast part only enqueues events/data and can signal to the
// source that it isn't ready for receive next portion. Slow part dequeues
// events/data and handles its. This scheme can be used to  implement interrupt 
// driven I/O library.
package main

import (
	"delay"

	"stm32/f4/exti"
	"stm32/f4/gpio"
	"stm32/f4/irq"
	"stm32/f4/periph"
	"stm32/f4/setup"
)

var LED = gpio.D

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

func init() {
	setup.Performance168(8)

	periph.AHB1ClockEnable(periph.GPIOA | periph.GPIOD)
	periph.AHB1Reset(periph.GPIOA | periph.GPIOD)

	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Orange, gpio.Out)
	LED.SetMode(Red, gpio.Out)
	LED.SetMode(Blue, gpio.Out)

	// Setup external interrupt source: user button.
	gpio.A.SetMode(0, gpio.In)
	exti.L0.Connect(gpio.A)
	exti.L0.RiseTrigEnable()
	exti.L0.IntEnable()
	irq.Ext0.UseHandler(buttonHandler)
	irq.Ext0.Enable()
}

func toggle(led) {
	LED.SetBit(led)
	delay.Millisec(500)
	LED.ClearBit(led)
	delay.Millisec(500)
}

var (
	c   = make(chan int, 3)
	led = Green
)

func buttonHandler() {
	exti.L0.ClearPending()
	// Non-blocking selects with buffered channels can be used in ISR.
	select {
	case c <- led:
		if led++; led > Red {
			led = Green
		}
	default:
		// Signal that c is full.
		LED.SetBit(Blue)
		delay.Loop(1e5)
		LED.ClearBit(Blue)
	}
}

func main() {
	for {
		toggle(<-c)
	}
}
