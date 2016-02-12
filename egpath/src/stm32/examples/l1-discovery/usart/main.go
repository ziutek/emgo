// This example shows how to use raw UART without interrupts (pulling).
// Used pins:
// - PB10: MCU Tx (connect to PC Rx),
// - PB11: MCU Rx (connect to PC Tx),
// - GND.
package main

import (
	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var tts *usart.USART

func init() {
	system.Setup32(0)
	systick.Setup()

	gpio.B.EnableClock(true)
	port, tx, rx := gpio.B, gpio.Pin10, gpio.Pin11

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART3)

	tts = usart.USART3

	tts.EnableClock(true)
	tts.SetBaudRate(115200)
	tts.SetConf(usart.RxEna | usart.TxEna)
	tts.Enable()
}

func writeByte(b byte) {
	tts.Store(int(b))
	for tts.Status()&usart.TxEmpty == 0 {
	}
}

func readByte() byte {
	for tts.Status()&usart.RxNotEmpty == 0 {
	}
	return byte(tts.Load())
}

func main() {
	for {
		b := readByte()
		writeByte(b)
		writeByte('\r')
		writeByte('\n')
	}
}
