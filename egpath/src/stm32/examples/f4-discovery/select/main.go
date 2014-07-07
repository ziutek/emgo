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

func toggle(led, d int) {
	LED.SetBit(led)
	delay.Loop(d)
	LED.ClearBit(led)
	delay.Loop(d)
}

const dly = 1e6

func blink(c <-chan int) {
	for {
		led := <-c
		toggle(led, dly*10)
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
			toggle(Orange, dly)
		case c2 <- Blue:
			toggle(Orange, dly)
		default:
			toggle(Green, dly)
		}
	}
}
