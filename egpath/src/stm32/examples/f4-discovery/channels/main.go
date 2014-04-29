package main

import (
	"delay"
	"runtime/noos"

	"stm32/f4/gpio"
	"stm32/f4/periph"
	"stm32/f4/setup"
)

var (
	LED = gpio.D
	In  = gpio.A
)

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

const Button = 0

func init() {
	setup.Performance(8)

	periph.AHB1ClockEnable(periph.GPIOA | periph.GPIOD)
	periph.AHB1Reset(periph.GPIOA | periph.GPIOD)

	In.SetMode(Button, gpio.In)

	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Orange, gpio.Out)
	LED.SetMode(Red, gpio.Out)
	LED.SetMode(Blue, gpio.Out)
}

func blink(c *Chan) {
	for {
		led := int(c.Recv())
		LED.SetBit(led)
		delay.Loop(1e7)
		LED.ClearBit(led)
		delay.Loop(1e7)
	}
}

func main() {
	c := Chan{event: noos.AssignEvent()}

	go blink(&c)
	go blink(&c)
	go blink(&c)

	led := Green
	for {
		c.Send(Elem(led))
		if led++; led > Blue {
			led = Green
		}
	}
}
