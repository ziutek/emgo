// This example shows how to use TIMER, GPIOTE and PPI to implement PWM output.
package main

import (
	"delay"
	"fmt"
	"rtos"

	"nrf5/ppipwm"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/irq"
	"nrf5/hal/ppi"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/timer"
	"nrf5/hal/uart"
)

var (
	pwm *ppipwm.Toggle
	u   *uart.Driver
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0

	u = uart.NewDriver(uart.UART0, make([]byte, 80))
	u.P.StorePSEL(uart.RXD, p0.Pin(11))
	u.P.StorePSEL(uart.TXD, p0.Pin(9))
	u.P.StoreBAUDRATE(uart.Baud115200)
	u.P.StoreENABLE(true)
	rtos.IRQ(irq.UART0).Enable()
	u.EnableTx()
	fmt.DefaultWriter = u

	pwm = ppipwm.NewToggle(timer.TIMER1)
	pwm.SetFreq(1, 4000)
	pwm.Setup(0, p0.Pin(19), gpiote.Chan(0), ppi.Chan(0), ppi.Chan(1))
	pwm.Setup(1, p0.Pin(20), gpiote.Chan(1), ppi.Chan(2), ppi.Chan(3))
	pwm.Setup(2, p0.Pin(21), gpiote.Chan(2), ppi.Chan(4), ppi.Chan(5))
}

func main() {
	max := pwm.Max()
	v, a, b, c := 0, 4, 5, 0
	for {
		fmt.Printf("%d/%d\r\n", v, max)
		pwm.SetVal(0, v)
		pwm.SetVal(1, v)
		pwm.SetVal(2, v)
		switch {
		case v == 0:
			c = a
			a, b = b, a
		case v > (max*5+4)/4:
			c = 0
			a, b = b, a
		}
		v = (v*a + c) / b
		delay.Millisec(200)
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
