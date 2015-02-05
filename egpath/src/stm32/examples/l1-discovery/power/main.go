// Control of power (water heater, house heating system).
package main

import (
	"rtos"

	"stm32/l1/exti"
	"stm32/l1/gpio"
	"stm32/l1/irqs"
	"stm32/l1/periph"
	"stm32/l1/setup"
	"stm32/l1/usarts"
	"stm32/usart"
)

var (
	ledsPort  = gpio.B
	waterPort = gpio.B
	waterExti = exti.L9
	ssrPort   = gpio.C
	onewPort  = gpio.C
	onewUART  = usarts.USART3
	onewClk   = setup.APB1Clk
)

const (
	blue  = LED(6)
	green = LED(7)
	water = uint(9)
	ssr0  = uint(6)
	ssr1  = uint(7)
	ssr2  = uint(8)
	onew  = uint(10)
)

func init() {
	setup.Performance(0)

	periph.APB1ClockEnable(periph.USART3)
	periph.APB1Reset(periph.USART3)
	periph.APB2ClockEnable(periph.SysCfg)
	periph.APB2Reset(periph.SysCfg)
	gpiop := ledsPort.Periph() | waterPort.Periph() | ssrPort.Periph() | onewPort.Periph()
	periph.AHBClockEnable(gpiop)
	periph.AHBReset(gpiop)

	// Setup LEDs output.

	ledsPort.SetMode(uint(green), gpio.Out)
	ledsPort.SetMode(uint(blue), gpio.Out)

	// Setup SSR output.

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

	// Setup USART to operate as 1-wire master.

	onewPort.SetMode(onew, gpio.Alt)
	onewPort.SetOutType(onew, gpio.OpenDrain)
	onewPort.SetAltFunc(onew, gpio.USART3)

	onewUART.SetWordLen(usart.Bits8)
	onewUART.SetParity(usart.None)
	onewUART.SetStopBits(usart.Stop1b)
	onewUART.SetMode(usart.Tx | usart.Rx)
	onewUART.SetHalfDuplex(true)
	onewUART.EnableIRQs(usart.RxNotEmptyIRQ)
	onewUART.Enable()

	rtos.IRQ(irqs.USART3).UseHandler(usart3__ISR)
	rtos.IRQ(irqs.USART3).Enable()

	periph.APB2ClockDisable(periph.SysCfg)
}

func ext9_5__ISR() {
	p := exti.Pending()
	(exti.L9 | exti.L8 | exti.L7 | exti.L6 | exti.L5).ClearPending()
	if waterExti&p != 0 {
		waterIRQ()
	}
}

func usart3__ISR() {
	onewSerial.IRQ()
}

func main() {
	waterTask()
}
