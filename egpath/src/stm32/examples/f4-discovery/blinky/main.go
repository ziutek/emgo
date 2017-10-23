package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var green, orange, red, blue gpio.Pin

func init() {
	system.Setup168(8)
	systick.Setup(2e6)

	gpio.D.EnableClock(false)
	green = gpio.D.Pin(12)
	orange = gpio.D.Pin(13)
	red = gpio.D.Pin(14)
	blue = gpio.D.Pin(15)

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}

	green.Setup(&cfg)
	orange.Setup(&cfg)
	red.Setup(&cfg)
	blue.Setup(&cfg)
}

func wait() {
	//delay.Loop(1e7)
	delay.Millisec(500)
}

func main() {
	for {
		green.Clear()
		orange.Set()
		wait()

		orange.Clear()
		red.Set()
		wait()

		red.Clear()
		blue.Set()
		wait()

		blue.Clear()
		green.Set()
		wait()
	}
}
