// This example shows how to use raw UART without interrupts (pulling).
// Used pins:
// - PB10: MCU Tx (connect to PC Rx),
// - PB11: MCU Rx (connect to PC Tx),
// - GND.
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
	system.Setup32(0)
	systick.Setup(2e6)

	gpio.B.EnableClock(true)
	port, tx, rx := gpio.B, gpio.Pin10, gpio.Pin11

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART3)

	tts = usart.USART3
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
