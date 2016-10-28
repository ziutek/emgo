package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var leds *gpio.Port

const (
	Green  = gpio.Pin12
	Orange = gpio.Pin13
	Red    = gpio.Pin14
	Blue   = gpio.Pin15
)

func init() {
	system.Setup168(8)
	systick.Setup()

	gpio.D.EnableClock(false)
	leds = gpio.D

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Orange|Red|Blue, &cfg)
}

func toggle(colors gpio.Pins) {
	leds.SetPins(colors)
	delay.Millisec(200)
	leds.ClearPins(colors)
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
