package main

import (
	"stm32/hal/gpio"
	"stm32/hal/setup"
	"stm32/hal/usart"
)

var tts *usart.USART

func init() {
	setup.Performance168(8)

	port, tx, rx := gpio.A, 2, 3

	port.EnableClock(true)
	port.SetMode(tx, gpio.Alt)
	port.SetAltFunc(tx, gpio.USART2)
	port.SetMode(rx, gpio.Alt)
	port.SetAltFunc(rx, gpio.USART2)

	tts = usart.USART2

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
