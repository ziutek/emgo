package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var led gpio.Pin

func init() {
	system.Setup96(8)
	systick.Setup(2e6)

	gpio.A.EnableClock(false)
	led = gpio.A.Pin(5)

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	led.Setup(&cfg)
}

func wait() {
	delay.Millisec(500)
}

func main() {
	for {
		led.Set()
		wait()
		led.Clear()
		wait()
	}
}
