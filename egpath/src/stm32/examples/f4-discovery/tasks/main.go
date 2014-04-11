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

func blink(l, d int, max, inc float32) {
	//for i := 0; i < 30 ; i++ {
	for inc < max {
		LED.SetBit(l)
		delay.Loop(d)
		LED.ClearBit(l)
		delay.Loop(d)
		inc *= inc
	}
}

func main() {
	for {
		go blink(Green, 8e5, 110, 1.0001)
		go blink(Orange, 5e5, 120, 1.0001)
		go blink(Red, 3e5, 130, 1.0001)
		blink(Blue, 12e5, 100, 1.0001)
		delay.Loop(1e7)
	}
}
