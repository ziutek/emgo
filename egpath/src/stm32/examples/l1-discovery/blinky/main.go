package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var leds *gpio.Port

const (
	Blue  = gpio.Pin6
	Green = gpio.Pin7
)

func init() {
	system.Setup32(0)
	systick.Setup()

	gpio.B.EnableClock(false)
	leds = gpio.B

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Blue, &cfg)
}

func main() {
	for {
		leds.ClearPins(Blue)
		leds.SetPins(Green)
		delay.Millisec(1000)

		leds.ClearPins(Green)
		leds.SetPins(Blue)
		delay.Millisec(1000)
	}
}
