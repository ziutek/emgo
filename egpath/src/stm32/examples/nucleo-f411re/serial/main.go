// This example shows how to use USART as serial console.
// It uses PA2, PA3 pins as Tx, Rx that are by default connected to ST-LINK
// Virtual Com Port.
package main

import (
	"delay"
	"fmt"
	"rtos"
	"strconv"

	"stm32/f4/gpio"
	"stm32/f4/irqs"
	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f4/usarts"
	"stm32/serial"
	"stm32/usart"
)

const (
	Green = 5
)

var (
	leds = gpio.A
	udev = usarts.USART2
	s    = serial.New(udev, 80, 8)
)

func init() {
	setup.Performance84(8)

	periph.AHB1ClockEnable(periph.GPIOA)
	periph.AHB1Reset(periph.GPIOA)

	leds.SetMode(Green, gpio.Out)

	periph.APB1ClockEnable(periph.USART2)
	periph.APB1Reset(periph.USART2)

	port, tx, rx := gpio.A, uint(2), uint(3)

	port.SetMode(tx, gpio.Alt)
	port.SetOutType(tx, gpio.PushPull)
	port.SetPull(tx, gpio.PullUp)
	port.SetOutSpeed(tx, gpio.Fast)
	port.SetAltFunc(tx, gpio.USART2)
	port.SetMode(rx, gpio.Alt)
	port.SetAltFunc(rx, gpio.USART2)

	udev.SetBaudRate(115200, setup.APB1Clk)
	udev.SetWordLen(usart.Bits8)
	udev.SetParity(usart.None)
	udev.SetStopBits(usart.Stop1b)
	udev.SetMode(usart.Tx | usart.Rx)
	udev.EnableIRQs(usart.RxNotEmptyIRQ)
	udev.Enable()

	rtos.IRQ(irqs.USART2).UseHandler(sirq)
	rtos.IRQ(irqs.USART2).Enable()

	s.SetUnix(true)
}

func blink(c uint, d int) {
	leds.SetBit(c)
	if d > 0 {
		delay.Millisec(d)
	} else {
		delay.Loop(-1e4 * d)
	}
	leds.ClearBit(c)
}

func sirq() {
	s.IRQ()
}

func checkErr(err error) {
	if err != nil {
		s.WriteString("\nError: ")
		s.WriteString(err.Error())
		s.WriteByte('\n')
	}
}

func main() {
	const (
		SmallestNormal         = 2.2250738585072014e-308
		SmallestNonzeroFloat64 = 4.940656458412465441765687928682213723651e-324
	)

	f := SmallestNonzeroFloat64

	wf, we, df, de := strconv.Show(f)

	fmt.Fprint(s, "w = ", wf, "p", we, "\n")
	fmt.Fprint(s, "d  = ", df, "p", de, "\n")

	var buf [50]byte
	n := strconv.FormatFloat64(buf[:], f)
	s.Write(buf[:n])
	s.WriteString("\n\n")
}
