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

type Chan struct {
	c *ChanA
	b []int
}

func NewChan(n int) Chan {
	return Chan{NewChanA(n), make([]int, n)}
}

func (c Chan) Send(e int) {
	n := c.c.Send()
	c.b[n] = e
	c.c.Done(n)
}

func (c Chan) Recv() int {
	n := c.c.Recv()
	e := c.b[n]
	c.c.Done(n)
	return e
}

func toggle(l, d int) {
	LED.SetBit(l)
	delay.Loop(d)
	LED.ClearBit(l)
	delay.Loop(d)
}

func blink(c Chan) {
	for {
		toggle(c.Recv(), 1e7)
	}
}

func main() {
	red := NewChan(6)
	green := NewChan(6)
	blue := NewChan(6)

	go blink(red)
	go blink(blue)
	go blink(green)

	for {
		red.Send(Red)
		toggle(Orange, 1e6)
		blue.Send(Blue)
		toggle(Orange, 1e6)
		green.Send(Green)
		toggle(Orange, 1e6)
	}
}
