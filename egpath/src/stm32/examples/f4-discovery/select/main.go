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

func toggle(led, dly int) {
	LED.SetPin(led)
	delay.Millisec(dly)
	LED.ClearPin(led)
	delay.Millisec(dly)
}

func blink(c <-chan int) {
	for {
		led := <-c
		toggle(led, 1200)
	}
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)

	// Consumers
	go blink(c1)
	go blink(c2)

	// Producer
	for {
		select {
		case c1 <- Red:
			toggle(Orange, 200)
		case c2 <- Blue:
			toggle(Orange, 200)
		default:
			toggle(Green, 200)
		}
	}
}
