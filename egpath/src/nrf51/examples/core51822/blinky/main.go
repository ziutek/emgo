package main

import (
	"rtos"

	"nrf51/hal/clock"
	"nrf51/hal/gpio"
	"nrf51/hal/irq"
	"nrf51/hal/rtc"
	"nrf51/hal/system"
	"nrf51/hal/system/timer/rtcst"
)

var (
	//emgo:const
	leds = [...]byte{18, 19, 20, 21, 22}

	p0 = gpio.P0
)

func init() {
	system.Setup(clock.Xtal, clock.Xtal, true)
	rtcst.Setup(rtc.RTC0, 1)

	for _, led := range leds {
		p0.SetMode(int(led), gpio.Out)
	}
}

func main() {
	led := leds[0]
	for i := 0; ; i++ {
		for rtos.Nanosec() < int64(i)*5e8 {
		}
		if i&1 != 0 {
			p0.SetPin(int(led))
		} else {
			p0.ClearPin(int(led))
		}
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
