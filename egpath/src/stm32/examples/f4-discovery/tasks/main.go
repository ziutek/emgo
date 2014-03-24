package main

import (
	"sync/barrier"

	"runtime/noos"

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

func main() {
	const wait = 2e6
	for {
		c := noos.Tick
		barrier.Compiler()

		if c%1373 == 0 {
			LED.SetBit(Green)
		} else {
			LED.ClearBit(Green)
		}
		if c%521 == 0 {
			LED.SetBit(Orange)
		} else {
			LED.ClearBit(Orange)
		}
		if c%251 == 0 {
			LED.SetBit(Red)
		} else {
			LED.ClearBit(Red)
		}
		if c%137 == 0 {
			LED.SetBit(Blue)
		} else {
			LED.ClearBit(Blue)
		}
	}
}
