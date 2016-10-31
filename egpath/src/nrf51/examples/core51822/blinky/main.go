package main

import (
	"delay"

	"nrf51/hal/clock"
	"nrf51/hal/gpio"
	"nrf51/hal/irq"
	"nrf51/hal/rtc"
	"nrf51/hal/system"
	"nrf51/hal/system/timer/rtcst"
)

const (
	led0 = 18 + iota
	led1
	led2
	led3
	led4
)

var p0 = gpio.P0

func init() {
	system.Setup(clock.Xtal, clock.Xtal, true)
	rtcst.Setup(rtc.RTC0, 1)

	for led := led0; led <= led4; led++ {
		p0.SetMode(led, gpio.Out)
	}
}

func main() {
	for {
		for led := led0; led <= led4; led++ {
			p0.ClearPins(1<<led0 | 1<<led1 | 1<<led2 | 1<<led3 | 1<<led4)
			p0.SetPin(led)
			delay.Millisec(500)
		}
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
