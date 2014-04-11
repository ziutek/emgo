package main

import (
	"delay"

	"stm32/f4/gpio"
	"stm32/f4/periph"
	"stm32/f4/setup"
)

var LED = gpio.D

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

	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Orange, gpio.Out)
	LED.SetMode(Red, gpio.Out)
	LED.SetMode(Blue, gpio.Out)
}

type Args struct {
	x, y, z int
	u, v, w int
}

func F(a Args, b int) int

func main() {
	a := Args{11, 22, 33, 44, 55, 66}
	wait := 2e6 + F(a, 77)
	wait += a.x
	for {
		LED.ClearBit(Green)
		LED.SetBit(Orange)
		delay.Loop(wait)

		LED.ClearBit(Orange)
		LED.SetBit(Red)
		delay.Loop(wait)

		LED.ClearBit(Red)
		LED.SetBit(Blue)
		delay.Loop(wait)

		LED.ClearBit(Blue)
		LED.SetBit(Green)
		delay.Loop(wait)
	}
}
