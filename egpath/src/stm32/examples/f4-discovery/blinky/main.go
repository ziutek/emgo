package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/setup"
)

var LED *gpio.Port

const (
	Green  = gpio.Pin12
	Orange = gpio.Pin13
	Red    = gpio.Pin14
	Blue   = gpio.Pin15
)

func init() {
	setup.Performance168(8)

	gpio.D.EnableClock(false)
	LED = gpio.D

	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	LED.Setup(Green|Orange|Red|Blue, cfg)
}

func wait() {
	//delay.Loop(1e7)
	delay.Millisec(500)
}

func main() {
	for {
		LED.ClearPins(Green)
		LED.SetPins(Orange)
		wait()

		LED.ClearPins(Orange)
		LED.SetPins(Red)
		wait()

		LED.ClearPins(Red)
		LED.SetPins(Blue)
		wait()

		LED.ClearPins(Blue)
		LED.SetPins(Green)
		wait()
	}
}
