package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var green, red gpio.Pin

func init() {
	system.Setup96(26)
	systick.Setup(2e6)

	gpio.B.EnableClock(false)
	green = gpio.B.Pin(0)
	red = gpio.B.Pin(1)

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}

	green.Setup(&cfg)
	red.Setup(&cfg)
}

func wait() {
	//delay.Loop(1e7)
	delay.Millisec(500)
}

func main() {
	for {
		green.Set()
		red.Set()
		wait()

		green.Clear()
		red.Clear()
		wait()
	}
}
