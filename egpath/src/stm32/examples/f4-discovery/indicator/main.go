// Touch some PD pins on your discovery.
package main

import (
	"stm32/f4/gpio"
	"stm32/f4/periph"
	"stm32/f4/setup"
)

var (
	LED = gpio.D
	In  = gpio.D
)

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

	In.SetMode(0, gpio.In)
	In.SetMode(1, gpio.In)
	In.SetMode(2, gpio.In)
	In.SetMode(3, gpio.In)
}

func main() {
	for {
		if In.Bit(0) {
			LED.SetBit(Red)
		} else {
			LED.ClearBit(Red)
		}
		if In.Bit(1) {
			LED.SetBit(Green)
		} else {
			LED.ClearBit(Green)
		}
		if In.Bit(2) {
			LED.SetBit(Blue)
		} else {
			LED.ClearBit(Blue)
		}
		if In.Bit(3) {
			LED.SetBit(Orange)
		} else {
			LED.ClearBit(Orange)
		}
	}
}
