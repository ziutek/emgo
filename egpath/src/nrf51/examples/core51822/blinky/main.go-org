package main

import (
	"delay"

	"nrf51/822/gpio"
	"nrf51/822/setup"
)

var leds = gpio.P0

const (
	Blue  = 6
	Green = 7
)

func init() {
	setup.Performance(0)

	//...

	leds.SetMode(Blue, gpio.Out)
	leds.SetMode(Green, gpio.Out)
}

func main() {
	for {
		leds.ClearBit(Blue)
		leds.SetBit(Green)
		delay.Millisec(1000)

		leds.ClearBit(Green)
		leds.SetBit(Blue)
		delay.Millisec(1000)
	}
}
