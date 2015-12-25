package main

import (
	"delay"
	
	"stm32/hal/gpio"
	"stm32/hal/setup"
)

var LED *gpio.Port

const (
	Green  = 12
	Orange = 13
	Red    = 14
	Blue   = 15
)

func init() {
	setup.Performance168(8)

	gpio.D.Enable(false)
	gpio.D.Reset()

	LED = gpio.D
	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Orange, gpio.Out)
	LED.SetMode(Red, gpio.Out)
	LED.SetMode(Blue, gpio.Out)
}

func wait() {
	//delay.Loop(1e7)
	delay.Millisec(500)
}

func main() {
	for {
		LED.ClearPin(Green)
		LED.SetPin(Orange)
		wait()

		LED.ClearPin(Orange)
		LED.SetPin(Red)
		wait()

		LED.ClearPin(Red)
		LED.SetPin(Blue)
		wait()

		LED.ClearPin(Blue)
		LED.SetPin(Green)
		wait()
	}
}
