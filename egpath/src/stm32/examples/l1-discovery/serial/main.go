// This example shows how to use USART as serial console.
//
// You need terminal emulator (eg. screen, minicom, hyperterm) and some USB to
// 3.3V TTL serial adapter (eg. FT232RL, CP2102 based). Warninig! If you use
// USB to RS232 it can destroy your Discovery board.

// Connct adapter's GND, Rx and Tx pins respectively to Discovery's GND, PB10,
// PB11 pins.
package main

import (
	"delay"
	"fmt"
	"runtime/noos"

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

func blink(c, d int) {
	leds.SetBit(c)
	if d > 0 {
		delay.Millisec(d)
	} else {
		delay.Loop(-2e3 * d)
	}
	leds.ClearBit(c)
}

func sirq() {
	//blink(Blue, -10) // Uncoment to see "hardware buffer overrun".
	s.IRQ()
}

func checkErr(err error) {
	if err != nil {
		blink(Blue, 10)
		s.WriteString("\nError: ")
		s.WriteString(err.Error())
		s.WriteByte('\n')
	}
}

func main() {
	s.WriteString("\nHello!\n")

	var uts [5]uint64
	for i := range uts {
		delay.Loop(2e3)
		uts[i] = noos.Uptime()
	}

	s.WriteString("\nFor loop:\n")

	for _, ut := range uts {
		fmt.Uint64(ut).Format(s, 10, -12)
		s.WriteString(" ns\n")
	}

	s.WriteString("Echo:\n")
	s.Flush()

	var buf [8]byte
	for {
		n, err := s.Read(buf[:])
		checkErr(err)

		ns := noos.Uptime()
		fmt.Uint64(ns).Format(s, 10, -12)

		s.WriteString(" ns \"")
		s.Write(buf[:n])
		s.WriteString("\"\n")

		blink(Green, 10)
	}
}
