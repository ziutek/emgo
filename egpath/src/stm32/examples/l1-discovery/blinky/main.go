package main

import (
	"delay"

	"stm32/l1/gpio"
	"stm32/l1/periph"
	"stm32/l1/setup"
)

var leds = gpio.B

const (
	Blue  = 6
	Green = 7
)

func init() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)

	leds.SetMode(Blue, gpio.Out)
	leds.SetMode(Green, gpio.Out)
}

func main() {
	for {
		leds.ClearBit(Blue)
		leds.SetBit(Green)
		delay.Millisec(1000)

		leds.ClearBit(Green)
		leds.SetBit(Blue)
		delay.Millisec(1000)
	}
}
