package main

import (
	"delay"
	"fmt"
	"rtos"
	"sync/fence"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/temp"
	"nrf5/hal/uart"
)

const bias = 7.00 * 4 // Specific for any chip (°C/4).

var (
	u    *uart.Driver
	done rtos.EventFlag
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	u = uart.NewDriver(uart.UART0, make([]byte, 80))
	u.P.StorePSEL(uart.SignalTXD, gpio.P0.Pin(9))
	u.P.StoreBAUDRATE(uart.Baud115200)
	u.P.StoreENABLE(true)
	rtos.IRQ(u.P.NVIC()).Enable()
	fmt.DefaultWriter = u
}

func main() {
	u.EnableTx()
	t := temp.TEMP
	t.Event(temp.DATARDY).Clear()
	t.Event(temp.DATARDY).EnableIRQ()
	rtos.IRQ(t.NVIC()).Enable()
	for {
		done.Reset(0)
		fence.W()
		t.Task(temp.START).Trigger()
		if done.Wait(1, rtos.Nanosec()+1e9) {
			f := float32(t.LoadTEMP()-bias) * 0.25
			fmt.Printf("T = %6.2f °C\r\n", f)
		} else {
			fmt.Printf("Timeout\r\n")
		}
		delay.Millisec(2e3)
	}
}

func uartISR() {
	u.ISR()
}

func tempISR() {
	temp.TEMP.Event(temp.DATARDY).Clear()
	done.Signal(1)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:  rtcst.ISR,
	irq.UART0: uartISR,
	irq.TEMP:  tempISR,
}
