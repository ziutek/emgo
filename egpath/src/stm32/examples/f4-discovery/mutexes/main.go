// This example was created before channels were implemented in Emgo. It shows
// example how to use mutexes to synchronize two gorutines.
package main

import (
	"delay"
	"sync"

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
	setup.Performance168(8)

	periph.AHB1ClockEnable(periph.GPIOA | periph.GPIOD)
	periph.AHB1Reset(periph.GPIOA | periph.GPIOD)

	In.SetMode(Button, gpio.In)

	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Orange, gpio.Out)
	LED.SetMode(Red, gpio.Out)
	LED.SetMode(Blue, gpio.Out)
}

func blink(led, d int) {
	for {
		LED.SetBit(led)
		delay.Millisec(d)
		LED.ClearBit(led)
		delay.Millisec(d)
	}
}

func toggle(m1, m2 *sync.Mutex) {
	leds := []int{Red, Orange, Blue}
	i := 0
	for {
		m1.Lock()
		m2.Unlock()
		LED.ClearBit(leds[i])
		i = (i + 1) % len(leds)
		LED.SetBit(leds[i])
		delay.Millisec(100)
	}
}

func main() {
	go blink(Green, 1000)

	var m1, m2 sync.Mutex
	m1.Lock()
	m2.Lock()
	go toggle(&m1, &m2)

	for {
		b := false
		for !b {
			b = In.Bit(Button)
		}
		m1.Unlock()
		m2.Lock()
	}
}
