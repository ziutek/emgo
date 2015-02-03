// This example shows how to use channels to divide interrupt handler into two
// parts: fast part - that runs in interrupt context and soft part - that runs
// in user context. Fast part only enqueues events/data and signals to the
// source if it isn't ready to receive next portion. Slow part dequeues
// events/data and handles them. This scheme can be used to implement
// interrupt driven I/O library.
package main

import (
	"delay"
	"rtos"

	"stm32/l1/exti"
	"stm32/l1/gpio"
	"stm32/l1/irqs"
	"stm32/l1/periph"
	"stm32/l1/setup"
)

var LED = gpio.B

const (
	Blue  = 6
	Green = 7
)

func init() {
	setup.Performance(0)

	periph.APB2ClockEnable(periph.SysCfg)
	periph.APB2Reset(periph.SysCfg)
	periph.AHBClockEnable(periph.GPIOA | periph.GPIOB)
	periph.AHBReset(periph.GPIOA | periph.GPIOB)

	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Blue, gpio.Out)

	// Setup external interrupt source: user button.
	gpio.A.SetMode(0, gpio.In)
	exti.L0.Connect(gpio.A)
	exti.L0.RiseTrigEnable()
	exti.L0.IntEnable()
	rtos.IRQ(irqs.Ext0).UseHandler(buttonHandler)
	rtos.IRQ(irqs.Ext0).Enable()

	periph.APB2ClockDisable(periph.SysCfg)
}

var c = make(chan uint, 3)

func buttonHandler() {
	exti.L0.ClearPending()
	select {
	case c <- Green:
	default:
		// Signal that c is full.
		LED.SetBit(Blue)
		delay.Loop(1e5)
		LED.ClearBit(Blue)
	}
}

func toggle(led uint) {
	LED.SetBit(led)
	delay.Millisec(500)
	LED.ClearBit(led)
	delay.Millisec(500)
}

func main() {
	for {
		toggle(<-c)
	}
}
