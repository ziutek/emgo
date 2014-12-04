// This example shows how to use USART as 1-wire master.
//
// Connct adapter's GND, Rx and Tx pins respectively to Discovery's GND, PD8,
// PD9 pins.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"stm32/f4/gpio"
	"stm32/f4/irqs"
	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f4/usarts"
	"stm32/onedrv"
	"stm32/serial"
	"stm32/usart"

	"onewire"
)

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

var (
	leds = gpio.D

	ow  = usarts.USART6
	one = serial.New(ow, 8, 8)

	us   = usarts.USART2
	term = serial.New(us, 80, 8)
)

func init() {
	setup.Performance168(8)

	// LEDS

	periph.AHB1ClockEnable(periph.GPIOD)
	periph.AHB1Reset(periph.GPIOD)

	leds.SetMode(Green, gpio.Out)
	leds.SetMode(Orange, gpio.Out)
	leds.SetMode(Red, gpio.Out)
	leds.SetMode(Blue, gpio.Out)

	// USART

	periph.AHB1ClockEnable(periph.GPIOA)
	periph.AHB1Reset(periph.GPIOA)
	periph.APB1ClockEnable(periph.USART2)
	periph.APB1Reset(periph.USART2)

	port, tx, rx := gpio.A, 2, 3

	port.SetMode(tx, gpio.Alt)
	port.SetOutType(tx, gpio.PushPull)
	port.SetPull(tx, gpio.PullUp)
	port.SetOutSpeed(tx, gpio.Fast)
	port.SetAltFunc(tx, gpio.USART2)
	port.SetMode(rx, gpio.Alt)
	port.SetAltFunc(rx, gpio.USART2)

	us.SetBaudRate(115200, setup.APB1Clk)
	us.SetWordLen(usart.Bits8)
	us.SetParity(usart.None)
	us.SetStopBits(usart.Stop1b)
	us.SetMode(usart.Tx | usart.Rx)
	us.EnableIRQs(usart.RxNotEmptyIRQ)
	us.Enable()

	rtos.IRQ(irqs.USART2).UseHandler(termISR)
	rtos.IRQ(irqs.USART2).Enable()

	term.SetUnix(true)

	// 1-wire

	periph.AHB1ClockEnable(periph.GPIOC)
	periph.AHB1Reset(periph.GPIOC)
	periph.APB2ClockEnable(periph.USART6)
	periph.APB2Reset(periph.USART6)

	port, tx = gpio.C, 6

	port.SetMode(tx, gpio.Alt)
	port.SetOutType(tx, gpio.OpenDrain)
	port.SetOutSpeed(tx, gpio.Fast)
	port.SetAltFunc(tx, gpio.USART6)

	ow.SetWordLen(usart.Bits8)
	ow.SetParity(usart.None)
	ow.SetStopBits(usart.Stop1b)
	ow.SetMode(usart.Tx | usart.Rx)
	ow.SetHalfDuplex(true)
	ow.EnableIRQs(usart.RxNotEmptyIRQ)
	ow.Enable()

	rtos.IRQ(irqs.USART6).UseHandler(oneISR)
	rtos.IRQ(irqs.USART6).Enable()
}

func termISR() {
	term.IRQ()
}

func oneISR() {
	one.IRQ()
}

func blink(c, d int) {
	leds.SetBit(c)
	if d > 0 {
		delay.Millisec(d)
	} else {
		delay.Loop(-1e4 * d)
	}
	leds.ClearBit(c)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	term.WriteString(err.Error())
	term.WriteByte('\n')
	for {
		blink(Red, 100)
		delay.Millisec(100)
	}
}

func ok() {
	term.WriteString("OK\n")
	blink(Green, 50)
}

func main() {
	drv := onedrv.UARTDriver{Serial: one, Clock: setup.APB2Clk}
	m := onewire.Master{Driver: &drv}

	term.WriteString("Searching for all devices on the bus...\n")
	s := onewire.MakeSearch(false)
	for m.SearchNext(&s) {
		fmt.Fprint(term, s.Dev(), fmt.N)
	}
	checkErr(s.Err())

	term.WriteString("Reading ROM code (valid if only one device connected)...\n")
	d, err := m.ReadROM()
	checkErr(err)

	fmt.Fprint(term, d, fmt.N)
}
