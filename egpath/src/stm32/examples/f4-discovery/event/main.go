// This example shows how to use rtos.EventFlag for communication between
// interrupt handler and thread.
package main

import (
	"rtos"

	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

const button = gpio.Pin0

var Red, Green, Blue, Orange gpio.Pin

func init() {
	system.Setup168(8)
	systick.Setup(2e6)

	gpio.A.EnableClock(true)
	btnport := gpio.A
	gpio.D.EnableClock(false)
	Green = gpio.D.Pin(12)
	Orange = gpio.D.Pin(13)
	Red = gpio.D.Pin(14)
	Blue = gpio.D.Pin(15)

	// LEDs

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	for _, pin := range []gpio.Pin{Green, Orange, Red, Blue} {
		pin.Setup(&cfg)
	}

	// Button

	btnport.Setup(button, &gpio.Config{Mode: gpio.In})
	line := exti.Lines(button)
	line.Connect(btnport)
	line.EnableRisiTrig()
	line.EnableIRQ()
	rtos.IRQ(irq.EXTI0).Enable()
}

var event rtos.EventFlag

func buttonISR() {
	exti.Lines(button).ClearPending()
	event.Signal(1)
}

func wait() {
	if event.Wait(1, rtos.Nanosec()+2e9) {
		event.Reset(0)
	}
}

func main() {
	for {
		Green.Clear()
		Orange.Set()
		wait()
		Orange.Clear()
		Red.Set()
		wait()
		Red.Clear()
		Blue.Set()
		wait()
		Blue.Clear()
		Green.Set()
		wait()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI0: buttonISR,
}
