// This example shows how to use USART as serial console.
// It uses PA2, PA3 pins as Tx, Rx that are by default connected to ST-LINK
// Virtual Com Port.
package main

import (
	"delay"
	"fmt"
	"reflect"
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
	var uts [10]uint64

	// Following loop disassembled:
	// 0:
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

	var i interface{}
	a := 4
	i = fmt.Int(a)
	switch v := i.(type) {
	case bool:
		s.WriteString("bool")
		if v {
			s.WriteString(" true")
		} else {
			s.WriteString(" false")
		}
	case int:
		s.WriteString("int")
		fmt.Int(v).Format(s, -20)
	case *int:
		s.WriteString("*int")
		fmt.Int(*v).Format(s, -20)
	case fmt.Int:
		s.WriteString("fmt.Int")
		v.Format(s, -10)
	default:
		s.WriteString("unk")
	}
	s.WriteString("\nReflection:\n")
	s.WriteString(reflect.TypeOf(i).String())
	s.WriteString("  kind: ")
	s.WriteString(reflect.TypeOf(i).Kind().String())
	s.WriteByte('\n')
	s.WriteString(reflect.TypeOf(a).String())
	s.WriteByte('\n')

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
