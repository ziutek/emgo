package main

import (
	"stm32/hal/gpio"
	"stm32/hal/osclk/systick"
	"stm32/hal/system"
	"stm32/hal/usart"
)

var tts *usart.USART

func init() {
	system.Setup168(8)
	systick.Setup()

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn})
	port.SetAltFunc(tx|rx, gpio.USART2)

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
