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

const dly = 1e6

func blink(color <-chan int, end chan<- struct{}) {
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
	color := make(chan int, 10)
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
