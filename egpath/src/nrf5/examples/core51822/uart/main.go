package main

import (
	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/uart"
)

var leds [5]gpio.Pin

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0
	for i := range leds {
		led := p0.Pin(18 + i)
		led.Setup(&gpio.Config{Mode: gpio.Out})
		leds[i] = led
	}

	u := uart.UART0
	u.SetPSEL(uart.SignalRXD, p0.Pin(11))
	u.SetPSEL(uart.SignalTXD, p0.Pin(9))
	u.SetBAUDRATE(0x00275000)
	u.SetENABLE(true)
}

func main() {
	u := uart.UART0

	u.Task(uart.STARTTX).Trigger()
	txdrdy := u.Event(uart.TXDRDY)

	s := "Hello!\r\n"
	for i := 0; ; i++ {
		leds[0].Store(i)
		for i := 0; i < len(s); i++ {
			txdrdy.Clear()
			u.SetTXD(s[i])
			for !txdrdy.IsSet() {
			}
		}
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
