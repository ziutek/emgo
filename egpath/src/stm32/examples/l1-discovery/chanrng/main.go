package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

func gen(c chan<- gpio.Pins, pins gpio.Pins) {
	for {
		c <- pins
	}
}

var leds = gpio.B

const (
	Blue  = gpio.Pin6
	Green = gpio.Pin7
)

func main() {
	system.Setup32(0)
	systick.Setup()

	gpio.B.EnableClock(false)
	leds = gpio.B

	leds.Setup(Green|Blue, &gpio.Config{Mode: gpio.Out, Speed: gpio.Low})

	cb := make(chan gpio.Pins, 2)
	cg := make(chan gpio.Pins, 2)

	go gen(cg, Green)
	go gen(cb, Blue)

	for {
		var led gpio.Pins
		select {
		case led = <-cg:
		case led = <-cb:
		}
		leds.SetPins(led)
		delay.Millisec(100)
		leds.ClearPins(led)
		delay.Millisec(100)
	}
}
