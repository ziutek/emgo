// This example shows how to use rtos.At function to implement deadline/timeout
// for communication with channels.
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
	Green  = gpio.Pin12
	Orange = gpio.Pin13
	Red    = gpio.Pin14
	Blue   = gpio.Pin15

	Button = gpio.Pin0
)

func init() {
	system.Setup168(8)
	systick.Setup(2e6)

	gpio.A.EnableClock(true)
	bport := gpio.A
	gpio.D.EnableClock(false)
	leds = gpio.D

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Orange|Red|Blue, &cfg)

	// Setup external interrupt source: user button.
	bport.Setup(Button, &gpio.Config{Mode: gpio.In})
	line := exti.Lines(Button)
	line.Connect(bport)
	line.EnableRisiTrig()
	line.EnableIRQ()

	rtos.IRQ(irq.EXTI0).Enable()
}

var (
	c   = make(chan struct{}, 3)
	led = Green
)

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

func toggle(pins gpio.Pins) {
	leds.SetPins(pins)
	delay.Millisec(300)
	leds.ClearPins(pins)
}

func main() {
	var colors = [...]gpio.Pins{Red, Green, Orange}
	i := 0
	for {
		select {
		case <-c:
		case <-rtos.At(rtos.Nanosec() + 2e9):
		}
		toggle(colors[i])
		if i++; i == 3 {
			i = 0
		}
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI0: buttonISR,
}
