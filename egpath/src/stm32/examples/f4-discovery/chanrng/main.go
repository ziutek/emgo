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

	gpio.D.Enable(false)

	LED = gpio.D
	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Orange, gpio.Out)
	LED.SetMode(Red, gpio.Out)
	LED.SetMode(Blue, gpio.Out)
}

func toggle(led int) {
	LED.SetPin(led)
	delay.Millisec(200)
	LED.ClearPin(led)
	delay.Millisec(200)
}

func gen(c chan<- struct{}) {
	for {
		c <- struct{}{}
	}
}

func main() {
	c0 := make(chan struct{})
	c1 := make(chan struct{})
	c2 := make(chan struct{}, 1)
	c3 := make(chan struct{}, 2)

	go gen(c0)
	go gen(c1)
	go gen(c2)
	go gen(c3)

	for {
		select {
		case <-c0:
			toggle(Red)
		case <-c1:
			toggle(Green)
		case <-c2:
			toggle(Blue)
		case <-c3:
			toggle(Orange)
		}
	}
}
