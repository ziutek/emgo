package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/setup"
)

var LED *gpio.Port

const (
	Green  = gpio.Pin12
	Orange = gpio.Pin13
	Red    = gpio.Pin14
	Blue   = gpio.Pin15
)

func init() {
	setup.Performance168(8)
	setup.UseSysTick()

	gpio.D.EnableClock(false)
	LED = gpio.D

	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	LED.Setup(Green|Orange|Red|Blue, cfg)
}

func toggle(leds gpio.Pins, dly int) {
	LED.SetPins(leds)
	delay.Millisec(dly)
	LED.ClearPins(leds)
	delay.Millisec(dly)
}

func blink(c <-chan gpio.Pins) {
	for {
		leds := <-c
		toggle(leds, 1200)
	}
}

func main() {
	c1 := make(chan gpio.Pins)
	c2 := make(chan gpio.Pins)

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
