// Control of power (water heater, house heating system).
package main

import (
	"rtos"

	"stm32/l1/exti"
	"stm32/l1/gpio"
	"stm32/l1/irq"
	"stm32/l1/periph"
	"stm32/l1/setup"
	"stm32/l1/usarts"
	"stm32/usart"
)

var (
	buttnPort = gpio.A
	buttnExti = exti.L0
	ledsPort  = gpio.B
	heatPort  = gpio.C
	waterPort = gpio.B
	waterExti = exti.L9
	ssrPort   = gpio.C
	onewPort  = gpio.C
	onewUART  = usarts.USART3
	onewClk   = setup.APB1Clk
)

const (
	// butnPort
	buttn = 0

	// ledsPort
	blue  = LED(6)
	green = LED(7)

	// heatPort
	heat0 = 0
	heat1 = 1
	heat2 = 2

	// waterPort
	water = 9

	// ssrPort
	ssr0 = 6
	ssr1 = 7
	ssr2 = 8

	// onewPort
	onew = 10
)

func init() {
	setup.Performance(0)

	periph.APB1ClockEnable(periph.USART3)
	periph.APB1Reset(periph.USART3)
	periph.APB2ClockEnable(periph.SysCfg)
	periph.APB2Reset(periph.SysCfg)
	gpiop := buttnPort.Periph() |
		ledsPort.Periph() |
		heatPort.Periph() |
		waterPort.Periph() |
		ssrPort.Periph() |
		onewPort.Periph()
	periph.AHBClockEnable(gpiop)
	periph.AHBReset(gpiop)

	// Setup SWO.
	gpio.B.SetMode(3, gpio.Alt)
	gpio.B.SetAltFunc(3, gpio.Sys)

	// Setup button input.
	buttnPort.SetMode(buttn, gpio.In)
	buttnExti.Connect(buttnPort)
	buttnExti.FallTrigEnable()
	buttnExti.IntEnable()
	rtos.IRQ(irq.Ext0).UseHandler(ext0__ISR)
	rtos.IRQ(irq.Ext0).Enable()

	// Setup LEDs output.

	ledsPort.SetMode(int(green), gpio.Out)
	ledsPort.SetMode(int(blue), gpio.Out)

	// Setup heating output.

	heatPort.SetMode(heat0, gpio.Out)
	heatPort.SetMode(heat1, gpio.Out)
	heatPort.SetMode(heat2, gpio.Out)

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
	rtos.IRQ(irq.Ext9_5).UseHandler(ext9_5__ISR)
	rtos.IRQ(irq.Ext9_5).Enable()

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

	rtos.IRQ(irq.USART3).UseHandler(usart3__ISR)
	rtos.IRQ(irq.USART3).Enable()

	periph.APB2ClockDisable(periph.SysCfg)
}

func ext0__ISR() {
	exti.L0.ClearPending()
	buttonIRQ()
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
	go heatingTask()
	waterTask()
}
