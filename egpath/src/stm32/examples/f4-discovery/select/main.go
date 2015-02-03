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

func toggle(led uint, d int) {
	LED.SetBit(led)
	delay.Millisec(d)
	LED.ClearBit(led)
	delay.Millisec(d)
}

func blink(c <-chan uint) {
	for {
		led := <-c
		toggle(led, 1200)
	}
}

func main() {
	c1 := make(chan uint)
	c2 := make(chan uint)

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
