package main

import (
	"rtos"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var tts *usart.Periph

func init() {
	system.Setup168(8)
	systick.Setup()

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3

	port.Setup(tx, gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART2)

	tts = usart.USART2
	tts.EnableClock(true)
	tts.SetBaudRate(115200)
	tts.Enable()
	tts.SetConf(usart.RxEna | usart.TxEna)
}

func writeByte(b byte) {
	for {
		ev, _ := tts.Status()
		if ev&usart.TxEmpty != 0 {
			break
		}
	}
	tts.Store(int(b))
}

func readByte() byte {
	for {
		ev, err := tts.Status()
		if err != 0 {
			tts.Load() // Reset errors.
			dbg := rtos.Debug(0)
			dbg.WriteString(err.Error())
			dbg.WriteByte('\n')
			continue
		}
		if ev&usart.RxNotEmpty != 0 {
			return byte(tts.Load())
		}
	}
}

func main() {
	for {
		b := readByte()
		writeByte(b)
		writeByte('\r')
		writeByte('\n')
	}
}
