package main

import (
	"fmt"
	"reflect"
	"rtos"

	"stm32/serial"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var con *serial.Dev

func init() {
	system.Setup96(8)
	systick.Setup()

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART2)

	s := usart.USART2

	s.EnableClock(true)
	s.SetBaudRate(115200)
	s.SetConf(usart.RxEna | usart.TxEna)
	s.EnableIRQs(usart.RxNotEmptyIRQ)
	s.Enable()

	con = serial.New(s, 80, 8)
	con.SetUnix(true)
	fmt.DefaultWriter = con

	rtos.IRQ(irq.USART2).Enable()
}

func conISR() {
	con.IRQ()
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

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2: conISR,
}
