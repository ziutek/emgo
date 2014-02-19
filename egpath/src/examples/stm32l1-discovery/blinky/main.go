package main

import (
	_ "cortexm/startup"
	"delay"
	"stm32/l1/gpio"
	"stm32/l1/periph"
)

// STM32L1-Discovery LEDs

var LEDs = gpio.B

const (
	Blue  = 6
	Green = 7
)

func init() {
	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)

	LEDs.SetMode(Blue, gpio.Out)
	LEDs.SetMode(Green, gpio.Out)
}

func main() {
	const wait = 2e4
	for {
		LEDs.ResetBit(Blue)
		LEDs.SetBit(Green)
		delay.Loop(wait)
		LEDs.ResetBit(Green)
		LEDs.SetBit(Blue)
		delay.Loop(wait)
	}
}
