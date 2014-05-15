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

func toggle(l, d int) {
	LED.SetBit(l)
	delay.Loop(d)
	LED.ClearBit(l)
	delay.Loop(d)
}

func blink(c <-chan int) {
	for {
		led := <-c
		toggle(led, 1e7)
	}
}

func main() {
	const n = 0 // Set n to 0, 1, 2, 4, ... and see LEDs.
	red := make(chan int, n)
	green := make(chan int, n)
	blue := make(chan int, n)

	// Consumers
	go blink(red)
	go blink(blue)
	go blink(green)

	// Producer
	for {
		red <- Red
		toggle(Orange, 1e6)
		blue <- Blue
		toggle(Orange, 1e6)
		green <- Green
		toggle(Orange, 1e6)
	}
}
