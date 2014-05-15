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

func toggle(led, d int) {
	LED.SetBit(led)
	delay.Loop(d)
	LED.ClearBit(led)
	delay.Loop(d)
}

func blink(c <-chan int) {
	for {
		led := <-c
		toggle(led, 1e7)
	}
}

func main() {
	c := make(chan int, 0)

	// Consumers
	go blink(c)
	go blink(c)
	go blink(c)

	// Producer
	for {
		c <- Red
		toggle(Orange, 1e6)
		c <- Blue
		toggle(Orange, 1e6)
		c <- Green
		toggle(Orange, 1e6)
	}
}
