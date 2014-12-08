// This example shows how to use USART as serial console.
//
// You need terminal emulator (eg. screen, minicom, hyperterm) and some USB to
// 3.3V TTL serial adapter (eg. FT232RL, CP2102 based). Warninig! If you use
// USB to RS232 it can destroy your Discovery board.

// Connct adapter's GND, Rx and Tx pins respectively to Discovery's GND, PA2,
// PA3 pins.
package main

import (
	"delay"
	"fmt"
	"rtos"

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

	port, tx, rx := gpio.A, 2, 3

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
	//blink(Blue, -10) // Uncoment to see "hardware buffer overrun".
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
	var uts [10]uint64

	// Following loop disassembled:
	// 0b:
	//   svc	5
	//   strd	r0, r1, [r3, #8]!
	//	 ldr	r2, [sp, #52]
	//   cmp	r3, r2
	//   bne.n  0b
	for i := range uts {
		uts[i] = rtos.Uptime()
	}

	s.WriteString("\nrtos.Uptime() in loop:\n")
	for i, ut := range uts {
		fmt.Int64(ut).Format(s, -12)
		s.WriteString(" ns")
		if i > 0 {
			fmt.Fprintf(s, " (dt = %v ns)", fmt.Int64(ut-uts[i-1]))
		}
		s.WriteByte('\n')
	}

	s.WriteString("Echo:\n")
	s.Flush()

	var buf [40]byte
	for {
		n, err := s.Read(buf[:])
		checkErr(err)

		ns := rtos.Uptime()
		fmt.Uint64(ns).Format(s, -12)

		s.WriteString(" ns \"")
		s.Write(buf[:n])
		s.WriteString("\"\n")

		blink(Green, 10)
	}
}
