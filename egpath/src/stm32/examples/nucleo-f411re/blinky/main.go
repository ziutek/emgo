package main

import (
	"delay"

	"stm32/f4/gpio"
	"stm32/f4/periph"
	"stm32/f4/setup"
)

var LED = gpio.A

const (
	Green = 5
)

func init() {
	setup.Performance84(8)

	periph.AHB1ClockEnable(periph.GPIOA)
	periph.AHB1Reset(periph.GPIOA)

	LED.SetMode(Green, gpio.Out)
}

func wait() {
	delay.Millisec(1000)
}

func main() {
	for {
		LED.SetPin(Green)
		wait()
		LED.ClearPin(Green)
		wait()
	}
}
