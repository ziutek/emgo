package main

import (
	"delay"
	//"fmt"

	"stm32/l1/gpio"
	"stm32/l1/irqs"
	"stm32/l1/periph"
	"stm32/l1/setup"
	"stm32/l1/usarts"
	"stm32/serial"
	"stm32/usart"
)

const (
	Blue  = 6
	Green = 7
)

var (
	leds = gpio.B
	udev = usarts.USART3
	s    = serial.NewSerial(udev, 80, 8)
)

func init() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)
	periph.APB1ClockEnable(periph.USART3)
	periph.APB1Reset(periph.USART3)

	leds.SetMode(Blue, gpio.Out)
	leds.SetMode(Green, gpio.Out)

	port, tx, rx := gpio.B, 10, 11

	port.SetMode(tx, gpio.Alt)
	port.SetOutType(tx, gpio.PushPullOut)
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
	udev.EnableIRQs(usart.RxNotEmptyIRQ)
	udev.Enable()

	irqs.USART3.UseHandler(sirq)
	irqs.USART3.Enable()

	s.SetUnix(true)
}

func blink(led int) {
	leds.SetBit(led)
	delay.Loop(1e5)
	leds.ClearBit(led)
}

func sirq() {
	//blink(Blue)
	s.IRQ()
}

func checkErr(err error) {
	if err != nil {
		blink(Blue)
		s.WriteString("\nError: ")
		s.WriteString(err.Error())
		s.WriteByte('\n')
	}
}

func main() {
	s.WriteString("Echo:\n")
	s.Flush()

	var buf [8]byte
	for {
		n, err := s.Read(buf[:])
		checkErr(err)

		s.Write(buf[:n])
		s.WriteByte('\n')

		blink(Green)
	}
}
