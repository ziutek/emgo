package main

import (
	"delay"
	"rtos"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/uart"
)

var (
	leds [5]gpio.Pin
	u    *uart.Driver
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0
	for i := range leds {
		led := p0.Pin(18 + i)
		led.Setup(gpio.Config{Mode: gpio.Out})
		leds[i] = led
	}

	u = uart.NewDriver(uart.UART0, make([]byte, 80))
	u.P.StorePSEL(uart.SignalRXD, p0.Pin(11))
	u.P.StorePSEL(uart.SignalTXD, p0.Pin(9))
	u.P.StoreBAUDRATE(uart.Baud115200)
	u.P.StoreENABLE(true)
	rtos.IRQ(u.P.IRQ()).Enable()
}

func main() {
	u.EnableRx()
	u.EnableTx()
	u.WriteString("\r\nHello World!\r\n")
	var buf [40]byte
	for i := 0; ; i++ {
		u.WriteByte('^')
		n, err := u.Read(buf[:])
		for i := 0; err != nil; i++ {
			leds[4].Store(i)
			delay.Millisec(200)
		}
		u.WriteByte('#')
		u.Write(buf[:n])
		leds[0].Store(i)
	}
}

func uartISR() {
	u.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:  rtcst.ISR,
	irq.UART0: uartISR,
}
