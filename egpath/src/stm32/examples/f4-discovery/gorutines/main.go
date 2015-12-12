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

func wait(ms int) {
	delay.Millisec(ms)
	//delay.Loop(ms * 1e4)
}

func blink(led int, d int, max, inc float32) {
	for inc < max {
		LED.SetPin(led)
		wait(d)
		LED.ClearPin(led)
		wait(d)
		// Use floating point calculations to test STMF4 FPU context switching.
		inc *= inc
	}
}

func main() {
	for {
		go blink(Green, 100, 110, 1.0001)
		go blink(Orange, 230, 120, 1.0001)
		go blink(Red, 350, 130, 1.0001)
		blink(Blue, 500, 100, 1.0001)
		wait(250)
		// BUG: In real application you schould ensure that all gorutines
		// finished before next loop. In this case Blue LED blinks longest
		// so this example works.
	}
}
