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
	buzzPort  = gpio.B
	waterExti = exti.L8
)

const (
	blue  = LED(6)
	green = LED(7)
	water = 8
	buzz  = 9
)

func init() {
	setup.Performance(0)

	periph.APB2ClockEnable(periph.SysCfg)
	periph.APB2Reset(periph.SysCfg)
	gpiop := ledsPort.Periph() | waterPort.Periph() | buzzPort.Periph()
	periph.AHBClockEnable(gpiop)
	periph.AHBReset(gpiop)

	// Setup LEDs
	ledsPort.SetMode(int(green), gpio.Out)
	ledsPort.SetMode(int(blue), gpio.Out)
	// Setup buzzer
	buzzPort.SetMode(buzz, gpio.Out)

	// Setup external interrupt source: water flow sensor.
	waterPort.SetMode(water, gpio.In)
	waterPort.SetPull(water, gpio.PullUp) // Noise prevention.
	waterExti.Connect(gpio.B)
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
