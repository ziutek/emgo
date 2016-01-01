package main

import (
	"fmt"
	"reflect"
	"rtos"

	"stm32/serial"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/setup"
	"stm32/hal/usart"
)

var con *serial.Dev

func init() {
	setup.Performance96(8)

	port, tx, rx := gpio.A, 2, 3

	port.EnableClock(true)
	port.SetMode(tx, gpio.Alt)
	port.SetAltFunc(tx, gpio.USART2)
	port.SetMode(rx, gpio.Alt)
	port.SetAltFunc(rx, gpio.USART2)

	tts := usart.USART2

	tts.EnableClock(true)
	tts.SetBaudRate(115200)
	tts.SetConf(usart.RxEna | usart.TxEna)
	tts.EnableIRQs(usart.RxNotEmptyIRQ)
	tts.Enable()

	con = serial.New(tts, 80, 8)
	con.SetUnix(true)
	fmt.DefaultWriter = con

	rtos.IRQ(irq.USART2).Enable()
}

func conISR() {
	con.IRQ()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2: conISR,
}

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
	fmt.Println("\nReflection test:\n")
	for _, iv := range ivals {
		v := reflect.ValueOf(iv)
		fmt.Println(v.Type())
	}
}
