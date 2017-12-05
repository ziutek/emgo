package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var led gpio.Pin

func init() {
	system.SetupPLL(8, 1, 72/8)
	systick.Setup(2e6)

	gpio.A.EnableClock(false)
	led = gpio.A.Pin(5)

	led.Setup(&gpio.Config{Mode: gpio.Out, Speed: gpio.Low})
}

func wait() {
	
}

func main() {
	for {
		led.Set()
		delay.Millisec(100)
		led.Clear()
		delay.Millisec(900)
	}
}
