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

func wait() {
	//delay.Loop(1e7)
	delay.Millisec(500)
}

func main() {
	for {
		LED.ClearBit(Green)
		LED.SetBit(Orange)
		wait()

		LED.ClearBit(Orange)
		LED.SetBit(Red)
		wait()

		LED.ClearBit(Red)
		LED.SetBit(Blue)
		wait()

		LED.ClearBit(Blue)
		LED.SetBit(Green)
		wait()
	}
}
