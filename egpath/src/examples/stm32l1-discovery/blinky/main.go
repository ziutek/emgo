package main

import (
	_ "cortexm/startup"
	"delay"
	"stm32/l1/gpio"
	"stm32/l1/periph"
	"stm32/l1/setup"
	"sync"
)

// STM32L1-Discovery LEDs

var LEDs = gpio.B

const (
	Blue  = 6
	Green = 7
)

func init() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)

	LEDs.SetMode(Blue, gpio.Out)
	LEDs.SetMode(Green, gpio.Out)
}

func main() {
	const wait = 1e6
	for {
		LEDs.ResetBit(Blue)
		LEDs.SetBit(Green)
		delay.Loop(wait)

		sync.Barrier()

		LEDs.ResetBit(Green)
		LEDs.SetBit(Blue)
		delay.Loop(wait)

		sync.Memory()
	}
}
