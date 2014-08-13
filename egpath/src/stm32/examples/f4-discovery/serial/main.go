package main

import (
	"stm32/f4/gpio"
	"stm32/f4/irq"
	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f4/usart"
	"stm32/serial"
)

var (
	udev = usart.USART2
	uirq = irq.USART2
)

func init() {
	setup.Performance168(8)

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

	uirq.UseHandler(sirq)
	uirq.Enable()

	udev.SetBaudRate(115200)
	udev.SetWordLen(usart.Bits8)
	udev.SetParity(usart.None)
	udev.SetStopBits(usart.Stop1b)
	udev.EnableIRQs(usart.TxEmptyIRQ)
	udev.Enable()
	udev.EnableTx()
	
}

var s = serial.NewSerial(udev)

func sirq() {
	s.IRQ()
}

func main() {
	
	for {
		for _, r := range []byte{'A', 'l', 'a', '!', '\r', '\n'} {
			s.WriteByte(r)
		}
	}
}
