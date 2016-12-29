package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var (
	ledport *gpio.Port
	ledpin  gpio.Pins
)

func init() {
	system.Setup(8, 1, 72/8)
	systick.Setup()

	gpio.A.EnableClock(false)
	ledport, ledpin = gpio.A, gpio.Pin5

	ledport.Setup(ledpin, &gpio.Config{Mode: gpio.Out, Speed: gpio.Low})
}

func wait() {
	delay.Millisec(500)
}

func main() {
	for {
		ledport.SetPins(ledpin)
		wait()
		ledport.ClearPins(ledpin)
		wait()
	}
}
