package main

import (
	"delay"
	"nrf51/gpio"
)

const (
	led0 = iota + 18
	led1
	led2
	led3
	led4
)

func main() {
	p := gpio.P0
	p.SetMode(led0, gpio.Out)
	p.SetMode(led1, gpio.Out)

	for {
		p.SetPins(1<<led0 | 1<<led1)
		delay.Loop(1e6)
		p.ClearPins(1<<led0 | 1<<led1)
		delay.Loop(1e6)
	}
}
