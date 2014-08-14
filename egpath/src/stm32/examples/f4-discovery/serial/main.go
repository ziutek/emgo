package main

import (
	"delay"

	"stm32/f4/gpio"
	"stm32/f4/irq"
	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f4/usart"
	"stm32/serial"
)

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

var (
	leds = gpio.D
	udev = usart.USART2
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

	udev.SetBaudRate(115200)
	udev.SetWordLen(usart.Bits8)
	udev.SetParity(usart.None)
	udev.SetStopBits(usart.Stop1b)
	udev.EnableIRQs(usart.RxNotEmptyIRQ)
	udev.SetMode(usart.Tx | usart.Rx)
	udev.Enable()

	irq.USART2.UseHandler(sirq)
	irq.USART2.Enable()
}

var s = serial.NewSerial(udev, 3, 3)

func blink(c, d int) {
	leds.SetBit(c)
	if d > 0 {
		delay.Millisec(d)
	} else {
		delay.Loop(-10e3 * d)
	}
	leds.ClearBit(c)
}

func sirq() {
	if s.IRQ() != nil {
		// Indicate Rx buffer overflow.
		blink(Red, -1)
	}
}

func main() {
	s.SetUnix(true)
	s.WriteString("Echo application\n\n")
	for {
		b, _ := s.ReadByte()
		s.WriteByte(b)
		blink(Green, 50) // 50 ms to see overflow LED blinking (red).
	}
}
