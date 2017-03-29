package main

import (
	"delay"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var (
	leds [5]gpio.Pin
	key  gpio.Pin
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0
	
	key = p0.Pin(16)
	for i := range leds {
		leds[i] = p0.Pin(18 + i)
	}

	for _, led := range leds {
		led.Setup(gpio.Config{Mode: gpio.Out})
	}
	key.Setup(gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
}

func main() {
	n := 0
	for {
		n += key.Load()*2 - 1
		switch n {
		case -1:
			n = len(leds) - 1
		case len(leds):
			n = 0
		}
		led := leds[n]
		led.Set()
		delay.Millisec(100)
		led.Clear()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
