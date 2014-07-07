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
	setup.Performance168(8)

	periph.AHB1ClockEnable(periph.GPIOD)
	periph.AHB1Reset(periph.GPIOD)

	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Orange, gpio.Out)
	LED.SetMode(Red, gpio.Out)
	LED.SetMode(Blue, gpio.Out)
}

func main() {
	for {
		LED.ClearBit(Green)
		LED.SetBit(Orange)
		delay.Millisec(500)

		LED.ClearBit(Orange)
		LED.SetBit(Red)
		delay.Millisec(500)

		LED.ClearBit(Red)
		LED.SetBit(Blue)
		delay.Millisec(500)

		LED.ClearBit(Blue)
		LED.SetBit(Green)
		delay.Millisec(500)
	}
}
