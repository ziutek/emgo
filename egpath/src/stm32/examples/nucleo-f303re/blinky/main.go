package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var led gpio.Pin

func init() {
	system.Setup(8, 1, 72/8)
	systick.Setup()

	gpio.A.EnableClock(false)
	led = gpio.A.Pin(5)

	led.Setup(&gpio.Config{Mode: gpio.Out, Speed: gpio.Low})
}

func wait() {
	delay.Millisec(250)
}

func main() {
	for {
		led.Set()
		wait()
		led.Clear()
		wait()
	}
}
