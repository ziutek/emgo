package main

import (
	"delay"

	"stm32/f4/exti"
	"stm32/f4/gpio"
	"stm32/f4/irqs"
	"stm32/f4/periph"
	"stm32/f4/setup"
)

var leds = gpio.D

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

func init() {
	setup.Performance168(8)

	periph.APB2ClockEnable(periph.SysCfg)
	periph.APB2Reset(periph.SysCfg)
	periph.AHB1ClockEnable(periph.GPIOC | periph.GPIOD)
	periph.AHB1Reset(periph.GPIOC | periph.GPIOD)

	leds.SetMode(Green, gpio.Out)
	leds.SetMode(Orange, gpio.Out)
	leds.SetMode(Red, gpio.Out)
	leds.SetMode(Blue, gpio.Out)

	gpio.C.SetMode(1, gpio.In)
	exti.L1.Connect(gpio.C)
	exti.L1.RiseTrigEnable()
	exti.L1.IntEnable()
	irqs.Ext1.UseHandler(pulse)
	irqs.Ext1.Enable()
}

func blink(led int) {
	leds.SetBit(led)
	delay.Loop(1e6)
	leds.ClearBit(led)
}

func pulse() {
	exti.L1.ClearPending()
	blink(Green)
}

func main() {
	for {
	}
}
