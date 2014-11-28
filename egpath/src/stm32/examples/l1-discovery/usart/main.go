// This example shows how to use raw UART without interrupts (pulling).
// Used pins:
// - PB10: MCU Tx (connect to PC Rx),
// - PB11: MCU Rx (connect to PC Tx),
// - GND.
package main

import (
	"stm32/l1/gpio"
	"stm32/l1/periph"
	"stm32/l1/setup"
	"stm32/l1/usarts"
	"stm32/usart"
)

var udev = usarts.USART3

func init() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)
	periph.APB1ClockEnable(periph.USART3)
	periph.APB1Reset(periph.USART3)

	port, tx, rx := gpio.B, 10, 11

	port.SetMode(tx, gpio.Alt)
	port.SetOutType(tx, gpio.PushPull)
	port.SetPull(tx, gpio.PullUp)
	port.SetOutSpeed(tx, gpio.Medium)
	port.SetAltFunc(tx, gpio.USART3)
	port.SetMode(rx, gpio.Alt)
	port.SetAltFunc(rx, gpio.USART3)

	udev.SetBaudRate(115200, setup.APB1Clk)
	udev.SetWordLen(usart.Bits8)
	udev.SetParity(usart.None)
	udev.SetStopBits(usart.Stop1b)
	udev.SetMode(usart.Tx | usart.Rx)
	udev.Enable()
}

func writeByte(b byte) {
	udev.Store(uint32(b))
	for udev.Status()&usart.TxEmpty == 0 {
	}
}

func readByte() byte {
	for udev.Status()&usart.RxNotEmpty == 0 {
	}
	return byte(udev.Load())
}

func main() {
	for {
		b := readByte()
		writeByte(b)
		writeByte('\r')
		writeByte('\n')
	}
}
