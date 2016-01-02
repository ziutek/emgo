package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/setup"
)

var LED *gpio.Port

const (
	Blue  = 6
	Green = 7
)

func init() {
	setup.Performance32(0)

	gpio.B.EnableClock(false)

	LED = gpio.B
	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Blue, gpio.Out)
}

func main() {
	for {
		LED.ClearPin(Blue)
		LED.SetPin(Green)
		delay.Millisec(1000)

		LED.ClearPin(Green)
		LED.SetPin(Blue)
		delay.Millisec(1000)
	}
}
