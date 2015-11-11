package main

import (
	"delay"

	"stm32/l1/gpio"
	"stm32/l1/periph"
	"stm32/l1/setup"
)

func gen(c chan<- int, v int) {
	for {
		c <- v
	}
}

var leds = gpio.B

const (
	Blue  = 6
	Green = 7
)

func main() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)

	leds.SetMode(Blue, gpio.Out)
	leds.SetMode(Green, gpio.Out)

	cb := make(chan int, 2)
	cg := make(chan int, 2)

	go gen(cg, Green)
	go gen(cb, Blue)

	for {
		var led int
		select {
		case led = <-cg:
		case led = <-cb:
		}
		leds.SetPin(led)
		delay.Millisec(100)
		leds.ClearPin(led)
		delay.Millisec(100)
	}
}
