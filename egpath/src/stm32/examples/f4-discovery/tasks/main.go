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

func blink(l, d int) {
	for {
		LED.SetBit(l)
		delay.Loop(d)
		LED.ClearBit(l)
		delay.Loop(d)
	}
}

func main() {
	go blink(Green, 1e6)
	go blink(Orange, 2e6)
	go blink(Red, 3e6)
	blink(Blue, 4e6)
}
