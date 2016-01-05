package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/setup"
)

type LED struct {
	Port *gpio.Port
	Pin  int
}

func (led LED) Init() {
	led.Port.SetMode(led.Pin, gpio.Out)
}

func (led LED) On() {
	led.Port.SetPin(led.Pin)
}

func (led LED) Off() {
	led.Port.ClearPin(led.Pin)
}

var leds = []LED{
	{gpio.B, 7},
	{gpio.B, 6},
	{gpio.B, 5},
	{gpio.D, 2},
}

func init() {
	setup.Performance(8, 72/8, false)

	gpio.B.EnableClock(true)
	gpio.D.EnableClock(true)

	for _, led := range leds {
		led.Init()
	}
}

func main() {
	i := 0
	d := -1
	n := len(leds) - 1
	for {
		leds[i].On()
		delay.Millisec(100)
		if i == 0 || i == n {
			d = -d
			delay.Millisec(50)
		}
		leds[i].Off()
		i += d
	}
}
