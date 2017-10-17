package main

import (
	"delay"
	"rtos"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var (
	leds [5]gpio.Pin
	key  gpio.Pin
)

const gch = gpiote.Chan(0)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0

	key = p0.Pin(16)
	for i := range leds {
		leds[i] = p0.Pin(18 + i)
	}

	for _, led := range leds {
		led.Setup(gpio.ModeOut)
	}
	key.Setup(gpio.ModeIn | gpio.PullUp)
	gch.Setup(key, gpiote.ModeEvent|gpiote.PolarityHiToLo|gpiote.OutInitHigh)
	gch.IN().Event().IRQ().Enable()
	rtos.IRQ(irq.GPIOTE).Enable()
}

func gpioteISR() {
	if ev := gch.IN().Event(); ev.IsSet() {
		ev.Clear()
		for _, led := range leds {
			led.Set()
		}
	}
}

func main() {
	gpioteISR()
	n := 0
	for {
		led := leds[n]
		led.Set()
		delay.Millisec(50)
		led.Clear()
		delay.Millisec(500)
		if n++; n >= len(leds) {
			n = 0
		}
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:   rtcst.ISR,
	irq.GPIOTE: gpioteISR,
}
