package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var led gpio.Pin

func init() {
	system.Setup()
	systick.Setup()

	gpio.A.EnableClock(false)
	led = gpio.A.Pin(5)

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	led.Setup(&cfg)
}

func main() {
	for {
		led.Set()
		delay.Millisec(50)
		led.Clear()
		delay.Millisec(950)
	}
}
