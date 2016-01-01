package main

import (
	"delay"

	"stm32/hal/gpio"
)

var leds = gpio.D

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

func initLEDs() {
	leds.EnableClock(false)
	leds.SetMode(Green, gpio.Out)
	leds.SetMode(Orange, gpio.Out)
	leds.SetMode(Red, gpio.Out)
	leds.SetMode(Blue, gpio.Out)
}

func blink(led, dly int) {
	leds.SetPin(led)
	if dly < 0 {
		delay.Loop(-dly * 1e3)
	} else {
		delay.Millisec(dly)
	}
	leds.ClearPin(led)
}
