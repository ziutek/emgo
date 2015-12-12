package main

import (
	"reflect"
	"rtos"

	"stm32/f4/gpio"
	"stm32/f4/irq"
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

	rtos.IRQ(irq.USART2).Enable()

	s.SetUnix(true)
}

func isr() {
	s.IRQ()
}

var IRQs = [...]func(){
	irq.USART2: isr,
} //c:__attribute__((section(".ISRs")))

type S struct {
	s string
}

func (s S) String() string {
	return s.s
}

type T struct {
	a int
	b byte
	S
}

type Stringer interface {
	String() string
}

var ivals = [...]interface{}{
	1,
	byte(2),
	uintptr(3),
	rtos.Debug(4),
	T{5, 6, S{"foo"}},
	&T{7, 8, S{"bar"}},
	nil,
}

func main() {
	s.WriteString("\nReflection test:\n\n")
	for _, iv := range ivals {
		v := reflect.ValueOf(iv)
		t := v.Type()
		s.WriteString(t.String())
		s.WriteByte('\n')
	}
}
