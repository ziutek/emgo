package main

import (
	_ "cortexm/startup"
	"delay"
	"stm32/f4/gpio"
	"stm32/f4/periph"
	"stm32/f4/setup"
)

var LEDs = gpio.D

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

func init() {
	setup.Performance(8)

	periph.AHB1ClockEnable(periph.GPIOD)
	periph.AHB1Reset(periph.GPIOD)

	LEDs.SetMode(Green, gpio.Out)
	LEDs.SetMode(Orange, gpio.Out)
	LEDs.SetMode(Red, gpio.Out)
	LEDs.SetMode(Blue, gpio.Out)
}

func main() {
	const wait = 1e6
	for {
		LEDs.ClearBit(Green)
		LEDs.SetBit(Orange)
		delay.Loop(wait)
		LEDs.ClearBit(Orange)
		LEDs.SetBit(Red)
		delay.Loop(wait)
		LEDs.ClearBit(Red)
		LEDs.SetBit(Blue)
		delay.Loop(wait)
		LEDs.ClearBit(Blue)
		LEDs.SetBit(Green)
		delay.Loop(wait)
	}
}
