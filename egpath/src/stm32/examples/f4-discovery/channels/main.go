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
	systick.Setup(2e6)

	gpio.D.EnableClock(false)
	leds = gpio.D

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Orange|Red|Blue, &cfg)
}

func toggle(colors gpio.Pins, d int) {
	leds.SetPins(colors)
	delay.Millisec(d)
	leds.ClearPins(colors)
	delay.Millisec(d)
}

const dly = 100

func blink(color <-chan gpio.Pins, end chan<- struct{}) {
	for {
		led, ok := <-color
		if !ok {
			end <- struct{}{}
			return
		}
		toggle(led, dly*10)
	}
}

func main() {
	color := make(chan gpio.Pins, 10)
	end := make(chan struct{}, 3)

	// Consumers
	go blink(color, end)
	go blink(color, end)
	go blink(color, end)

	// Producer
	for i := 0; i < 10; i++ {
		color <- Red
		toggle(Orange, dly)
		color <- Blue
		toggle(Orange, dly)
		color <- Green
		toggle(Orange, dly)
	}
	close(color)

	// Wait for consumers.
	<-end
	<-end
	<-end
	for {
		toggle(Orange, dly)
	}
}
