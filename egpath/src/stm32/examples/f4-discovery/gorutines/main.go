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
	leds.Setup(Green|Orange|Red|Blue, &cfg)
}

func wait(ms int) {
	delay.Millisec(ms)
	//delay.Loop(ms * 1e4)
}

func blink(colors gpio.Pins, d int, max, inc float32) {
	for inc < max {
		leds.SetPins(colors)
		wait(d)
		leds.ClearPins(colors)
		wait(d)
		// Use floating point calculations to test STMF4 FPU context switching.
		inc *= inc
	}
}

func main() {
	for {
		go blink(Green, 100, 110, 1.0001)
		go blink(Orange, 230, 120, 1.0001)
		go blink(Red, 350, 130, 1.0001)
		blink(Blue, 500, 100, 1.0001)
		wait(250)
		// BUG: In real application you schould ensure that all gorutines
		// finished before next loop. In this case Blue LED blinks longest
		// so this example works.
	}
}
