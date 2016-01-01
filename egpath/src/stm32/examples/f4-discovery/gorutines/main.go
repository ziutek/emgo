package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/setup"
)

var LED *gpio.Port

const (
	Green  = 12
	Orange = 13
	Red    = 14
	Blue   = 15
)

func init() {
	setup.Performance168(8)

	gpio.D.EnableClock(false)

	LED = gpio.D
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
