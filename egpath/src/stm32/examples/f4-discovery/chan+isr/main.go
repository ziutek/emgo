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
	setup.Performance(8)

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

func toggle(led, d int) {
	LED.SetBit(led)
	delay.Loop(d)
	LED.ClearBit(led)
	delay.Loop(d)
}

const dly = 1e6

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
		toggle(Blue, dly)
	}
}

func main() {
	for {
		toggle(<-c, 5*dly)
	}
}
