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
		LED.Clear(Green)
		LED.Set(Orange)
		wait()

		LED.Clear(Orange)
		LED.Set(Red)
		wait()

		LED.Clear(Red)
		LED.Set(Blue)
		wait()

		LED.Clear(Blue)
		LED.Set(Green)
		wait()
	}
}
