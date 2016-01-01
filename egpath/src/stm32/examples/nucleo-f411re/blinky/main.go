package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/setup"
)

var LEDport = gpio.A

const Green = 5

func init() {
	setup.Performance96(8)
	LEDport.EnableClock(false)
	LEDport.SetMode(Green, gpio.Out)
}

func wait() {
	delay.Millisec(500)
}

func main() {
	for {
		LEDport.SetPin(Green)
		wait()
		LEDport.ClearPin(Green)
		wait()
	}
}
