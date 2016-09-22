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
	systick.Setup()

	gpio.D.EnableClock(false)
	leds = gpio.D

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Orange|Red|Blue, cfg)
}

func wait() {
	//delay.Loop(1e7)
	delay.Millisec(500)
}

func main() {
	for {
		leds.ClearPins(Green)
		leds.SetPins(Orange)
		wait()

		leds.ClearPins(Orange)
		leds.SetPins(Red)
		wait()

		leds.ClearPins(Red)
		leds.SetPins(Blue)
		wait()

		leds.ClearPins(Blue)
		leds.SetPins(Green)
		wait()
	}
}
