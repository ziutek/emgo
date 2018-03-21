// This example demonstrates some feutures of Go language: gorutines, channels,
// empty interface and type switch. This program looks like Go program and works
// like Go program but can run on MCU that has only 4 KB RAM and 16 KB Flash.
package main

import (
	"delay"
	"math/rand"
	"rtos"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var led gpio.Pin

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	gpio.A.EnableClock(false)
	led = gpio.A.Pin(4)

	cfg := &gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low}
	led.Setup(cfg)
}

type On int

func task1(ch chan<- interface{}) {
	var r rand.XorShift64
	r.Seed(rtos.Nanosec())
	for {
		ch <- On(rnd.Int32()&511 + 89)
	}
}

type Off int

func task2(ch chan<- interface{}) {
	var rnd rand.XorShift64
	rnd.Seed(rtos.Nanosec())
	for {
		ch <- Off(rnd.Int32()&511 + 50)
	}
}

func main() {
	ch := make(chan interface{})
	go task1(ch)
	go task2(ch)
	for val := range ch {
		switch ms := val.(type) {
		case On:
			led.Clear() // Turn LED on.
			delay.Millisec(int(ms))
		case Off:
			led.Set() // Turn LED off.
			delay.Millisec(int(ms))
		}
	}
}
