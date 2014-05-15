// This example is for people that are interested in channels implementation.
// See ../channels for example of channels usage.
package main

import (
	"delay"
	"unsafe"

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
	a *ChanA
	s *ChanS
}

func NewChan(n int) Chan {
	if n == 0 {
		return Chan{s: NewChanS()}
	}
	return Chan{a: NewChanA(n, unsafe.Sizeof(int(0)), unsafe.Alignof(int(0)))}
}

func (c Chan) Send(e int) {
	send, done := c.a.Send, c.a.Done
	if c.s != nil {
		send, done = c.s.Send, c.s.Done
	}
	p, d := send(unsafe.Pointer(&e))
	if p == nil {
		return
	}
	*(*int)(p) = e
	done(d)
}

func (c Chan) Recv() (e int, ok bool) {
	recv, done := c.a.Recv, c.a.Done
	if c.s != nil {
		recv, done = c.s.Recv, c.s.Done
	}
	p, d := recv(unsafe.Pointer(&e))
	if p == nil {
		return e, (d == 0)
	}
	e = *(*int)(p)
	done(d)
	return e, true
}

func toggle(l, d int) {
	LED.SetBit(l)
	delay.Loop(d)
	LED.ClearBit(l)
	delay.Loop(d)
}

func blink(c Chan) {
	for {
		led, _ := c.Recv()
		toggle(led, 1e7)
	}
}

func main() {
	c := NewChan(10)

	// Consumers
	go blink(c)
	go blink(c)
	go blink(c)

	// Producer
	for {
		c.Send(Red)
		toggle(Orange, 1e6)
		c.Send(Blue)
		toggle(Orange, 1e6)
		c.Send(Green)
		toggle(Orange, 1e6)
	}
}
