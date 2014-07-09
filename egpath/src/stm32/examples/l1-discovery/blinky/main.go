package main

import (
	"delay"

	"stm32/l1/gpio"
	"stm32/l1/periph"
	"stm32/l1/setup"
)

var LED = gpio.B

const (
	Blue  = 6
	Green = 7
)

func init() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)

	LED.SetMode(Blue, gpio.Out)
	LED.SetMode(Green, gpio.Out)
}

func main() {
	for {
		LED.ClearBit(Blue)
		LED.SetBit(Green)
		delay.Millisec(1000)

		LED.ClearBit(Green)
		LED.SetBit(Blue)
		delay.Millisec(1000)
	}
}
