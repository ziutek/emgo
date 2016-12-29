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

var leds struct{ Red, Green, Blue, Orange gpio.Pin }

func init() {
	system.Setup168(8)
	systick.Setup()

	gpio.A.EnableClock(true)
	btnport := gpio.A
	gpio.D.EnableClock(false)
	leds.Green = gpio.D.Pin(12)
	leds.Orange = gpio.D.Pin(13)
	leds.Red = gpio.D.Pin(14)
	leds.Blue = gpio.D.Pin(15)

	// LEDs

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	pins := []gpio.Pin{leds.Green, leds.Orange, leds.Red, leds.Blue}
	for _, pin := range pins {
		pin.Port().SetupPin(pin.Index(), &cfg)
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
		leds.Green.Clear()
		leds.Orange.Set()
		wait()
		leds.Orange.Clear()
		leds.Red.Set()
		wait()
		leds.Red.Clear()
		leds.Blue.Set()
		wait()
		leds.Blue.Clear()
		leds.Green.Set()
		wait()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI0: buttonISR,
}
