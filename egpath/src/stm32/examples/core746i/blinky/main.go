// This example blinks two LEDs that need to be connected to PH3, PH5 pins. It
// does not use LEDs on the mother board (i have no one). If you want to use
// LEDs on the mother board change line:
//
//	leds, led1, led2 = gpio.H, gpio.Pin3, gpio.Pin4
//
// to
//
//  leds, led1, led2 = gpio.B, gpio.Pin6, gpio.Pin7
//
package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var (
	leds       *gpio.Port
	led1, led2 gpio.Pins
)

func init() {
	system.Setup192(8)
	systick.Setup()

	gpio.H.EnableClock(false)
	leds, led1, led2 = gpio.H, gpio.Pin3, gpio.Pin4

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(led1|led2, &cfg)
}

func wait() {
	//delay.Loop(1e7)
	delay.Millisec(500)
}

func main() {
	for {
		leds.ClearPins(led1)
		leds.SetPins(led2)
		wait()

		leds.ClearPins(led2)
		leds.SetPins(led1)
		wait()
	}
}
