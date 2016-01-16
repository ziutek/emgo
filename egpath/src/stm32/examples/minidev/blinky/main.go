package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/setup"
)

var leds *gpio.Port

const Blue = gpio.Pin13

func init() {
	setup.Performance(8, 72/8, false)

	gpio.C.EnableClock(true)
	leds = gpio.C

	leds.Setup(Blue, &gpio.Config{Mode: gpio.Out, Speed: gpio.Low})
}

func wait() {
	//delay.Loop(1e7)
	delay.Millisec(100)
}

func main() {
	for {
		leds.SetPins(Blue)
		wait()
		leds.ClearPins(Blue)
		wait()
	}
}
