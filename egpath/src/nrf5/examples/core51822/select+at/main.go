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
	ch   = make(chan struct{}, 1)
)

const key = gpiote.Chan(0)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0

	keyPin := p0.Pin(16)
	for i := range leds {
		leds[i] = p0.Pin(18 + i)
	}

	for _, led := range leds {
		led.Setup(gpio.ModeOut)
	}
	keyPin.Setup(gpio.ModeIn | gpio.PullUp)
	key.Setup(keyPin, gpiote.ModeEvent|gpiote.PolarityHiToLo)
	key.IN().Event().EnableIRQ()
	rtos.IRQ(irq.GPIOTE).Enable()
}

func main() {
	n := 0
	for {
		leds[n].Set()
		delay.Millisec(50)
		leds[n].Clear()
		if n++; n >= len(leds) {
			n = 0
		}
		select {
		case <-ch:
		case <-rtos.At(rtos.Nanosec() + 2e9):
		}
	}
}

func gpioteISR() {
	if ev := key.IN().Event(); ev.IsSet() {
		ev.Clear()
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:   rtcst.ISR,
	irq.GPIOTE: gpioteISR,
}
