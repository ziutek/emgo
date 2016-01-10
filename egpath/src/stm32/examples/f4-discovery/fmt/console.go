package main

import (
	"fmt"
	"rtos"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/usart"
	"stm32/serial"
)

var con *serial.Dev

func initConsole() {
	gpio.A.EnableClock(true)
	port, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3

	port.Setup(tx, &gpio.Config{Mode: gpio.Alt})
	port.Setup(rx, &gpio.Config{Mode: gpio.AltIn})
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

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2: conISR,
}
