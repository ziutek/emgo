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
	"rtos"

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
	s    = serial.New(udev, 80, 8)
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
	port.SetOutType(tx, gpio.PushPull)
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

	rtos.IRQ(irqs.USART3).UseHandler(sirq)
	rtos.IRQ(irqs.USART3).Enable()

	s.SetUnix(true)
}

func blink(c, dly int) {
	leds.SetPin(c)
	if dly > 0 {
		delay.Millisec(dly)
	} else {
		delay.Loop(-2e3 * dly)
	}
	leds.ClearPin(c)
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
		fmt.Fprint(s, ut, " ns")
		if i > 0 {
			fmt.Fprint(s, " (dt = ", ut-uts[i-1], " ns)")
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
		fmt.Fprint(s, ns, " ns \"")
		s.Write(buf[:n])
		s.WriteString("\"\n")

		blink(Green, 10)
	}
}
