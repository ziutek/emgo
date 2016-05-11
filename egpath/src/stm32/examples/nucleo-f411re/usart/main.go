// This example shows how to use USART as serial console.
// It uses PA2, PA3 pins as Tx, Rx that are by default connected to ST-LINK
// Virtual Com Port.
package main

import (
	"fmt"
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

	// GPIO

	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3

	// USART

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	port.SetAltFunc(tx|rx, gpio.USART2)

	s := usart.USART2

	s.EnableClock(true)
	s.SetBaudRate(115200)
	s.SetConf(usart.RxEna | usart.TxEna)
	s.EnableIRQ(usart.RxNotEmptyIRQ)
	s.Enable()

	con = serial.New(s, 80, 8)
	con.SetUnix(true)
	fmt.DefaultWriter = con

	rtos.IRQ(irq.USART2).Enable()
}

func conISR() {
	con.IRQ()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2: conISR,
}

func main() {
	fmt.Println("USART test:")
	buf := make([]byte, 80)
	for i := 0; ; i++ {
		n, err := con.Read(buf)
		fmt.Printf("%d %v '%s'\n", i, err, buf[:n])
	}
}
