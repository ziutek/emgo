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

func blink(led, d int, max, inc float32) {
	for inc < max {
		LED.SetBit(led)
		delay.Loop(d)
		LED.ClearBit(led)
		delay.Loop(d)
		inc *= inc
	}
}

func main() {
	for {
		go blink(Green, 11e5, 110, 1.0001)
		go blink(Orange, 7e5, 120, 1.0001)
		go blink(Red, 3e5, 130, 1.0001)
		blink(Blue, 17e5, 100, 1.0001)
		delay.Loop(1e7)
		// BUG: In real application you schould ensure that all gorutines
		// finished before next loop. In this case Blue LED blinks longest
		// so this works.
	}
}
