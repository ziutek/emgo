// This example shows how to use USART as serial console.
//
// You need terminal emulator (eg. screen, minicom, hyperterm) and some USB to
// 3.3V TTL serial adapter (eg. FT232RL, CP2102 based). Warninig! If you use
// USB to RS232 it can destroy your Discovery board.

// Connct adapter's GND, Rx and TX pins respectively to Discovery's GND, PA2,
// PA3 pins.
package main

import (
	"delay"

	"stm32/f4/gpio"
	"stm32/f4/irqs"
	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f4/usarts"
	"stm32/serial"
	"stm32/usart"
)

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

var (
	leds = gpio.D
	udev = usarts.USART2
	s    = serial.NewSerial(udev, 80, 8)
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

	io, tx, rx := gpio.A, 2, 3

	io.SetMode(tx, gpio.Alt)
	io.SetOutType(tx, gpio.PushPullOut)
	io.SetPull(tx, gpio.PullUp)
	io.SetOutSpeed(tx, gpio.Fast)
	io.SetAltFunc(tx, gpio.USART2)
	io.SetMode(rx, gpio.Alt)
	io.SetAltFunc(rx, gpio.USART2)

	udev.SetBaudRate(115200, setup.APB1Clk)
	udev.SetWordLen(usart.Bits8)
	udev.SetParity(usart.None)
	udev.SetStopBits(usart.Stop1b)
	udev.SetMode(usart.Tx | usart.Rx)
	udev.EnableIRQs(usart.RxNotEmptyIRQ)
	udev.Enable()

	irqs.USART2.UseHandler(sirq)
	irqs.USART2.Enable()

	s.SetUnix(true)
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

func sirq() {
	// blink(Blue, -10) // Uncoment to see "hardware buffer overrun" error.
	s.IRQ()
}

func checkErr(err error) {
	if err != nil {
		blink(Red, 10)
		s.WriteString("\nError: ")
		s.WriteString(err.Error())
		s.WriteByte('\n')
	}
}

func main() {
	s.WriteString("Echo application\n\n")
	s.Flush()
	var buf [40]byte
	for {
		n, err := s.Read(buf[:])
		checkErr(err)
		s.Write(buf[:n])
		blink(Green, 10)
	}
}
