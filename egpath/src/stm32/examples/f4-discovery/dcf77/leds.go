package main

import (
	"delay"

	"stm32/hal/gpio"
)

var leds *gpio.Port

const (
	Green  = gpio.Pin12
	Orange = gpio.Pin13
	Red    = gpio.Pin14
	Blue   = gpio.Pin15
)

func initLEDs(port *gpio.Port) {
	leds = port
	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Orange|Red|Blue, cfg)
}

func blink(pins gpio.Pins, dly int) {
	leds.SetPins(pins)
	if dly < 0 {
		delay.Loop(-dly * 1e3)
	} else {
		delay.Millisec(dly)
	}
	leds.ClearPins(pins)
}
