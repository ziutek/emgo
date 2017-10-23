package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var leds *gpio.Port

const (
	Green  = gpio.Pin12
	Orange = gpio.Pin13
	Red    = gpio.Pin14
	Blue   = gpio.Pin15
)

func init() {
	system.Setup168(8)
	systick.Setup(2e6)

	gpio.D.EnableClock(false)
	leds = gpio.D

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Orange|Red|Blue, &cfg)
}

func toggle(colors gpio.Pins, dly int) {
	leds.SetPins(colors)
	delay.Millisec(dly)
	leds.ClearPins(colors)
	delay.Millisec(dly)
}

func blink(c <-chan gpio.Pins) {
	for {
		colors := <-c
		toggle(colors, 1200)
	}
}

func main() {
	c1 := make(chan gpio.Pins)
	c2 := make(chan gpio.Pins)

	// Consumers
	go blink(c1)
	go blink(c2)

	// Producer
	for {
		select {
		case c1 <- Red:
			toggle(Orange, 200)
		case c2 <- Blue:
			toggle(Orange, 200)
		default:
			toggle(Green, 200)
		}
	}
}
