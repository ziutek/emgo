// Control of power (water heater, house heating system).
package main

import (
	"rtos"

	"stm32/l1/exti"
	"stm32/l1/gpio"
	"stm32/l1/irqs"
	"stm32/l1/periph"
	"stm32/l1/setup"
)

var (
	ledsPort  = gpio.B
	waterPort = gpio.B
	waterExti = exti.L9
	ssrPort   = gpio.C
)

const (
	blue  = LED(6)
	green = LED(7)
	water = 9
	ssr0  = 6
	ssr1  = 7
	ssr2  = 8
)

func init() {
	setup.Performance(0)

	periph.APB2ClockEnable(periph.SysCfg)
	periph.APB2Reset(periph.SysCfg)
	gpiop := ledsPort.Periph() | waterPort.Periph() | ssrPort.Periph()
	periph.AHBClockEnable(gpiop)
	periph.AHBReset(gpiop)

	// Setup LEDs output.
	ledsPort.SetMode(int(green), gpio.Out)
	ledsPort.SetMode(int(blue), gpio.Out)
	// Setup SSR output
	ssrPort.SetMode(ssr0, gpio.Out)
	ssrPort.SetMode(ssr1, gpio.Out)
	ssrPort.SetMode(ssr2, gpio.Out)

	// Setup external interrupt source: water flow sensor.
	waterPort.SetMode(water, gpio.In)
	waterPort.SetPull(water, gpio.PullUp) // Noise prevention.
	waterExti.Connect(waterPort)
	waterExti.FallTrigEnable()
	waterExti.IntEnable()
	rtos.IRQ(irqs.Ext9_5).UseHandler(ext9_5__ISR)
	rtos.IRQ(irqs.Ext9_5).Enable()

	periph.APB2ClockDisable(periph.SysCfg)
}

func ext9_5__ISR() {
	p := exti.Pending()
	(exti.L9 | exti.L8 | exti.L7 | exti.L6 | exti.L5).ClearPending()
	if waterExti&p != 0 {
		waterIRQ()
	}
}

func main() {
	waterTask()
}
