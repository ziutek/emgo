// This example shows how to use GPIO to as digital output (to blink connected
// LEDs) and input (to read state of KEY1).
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
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// Allocate pins (always do it in one place to avoid conflicts).
	p0 := gpio.P0
	key = p0.Pin(16) // KEY1
	for i := range leds {
		leds[i] = p0.Pin(18 + i)
	}

	// Configure pins.
	for _, led := range leds {
		led.Setup(gpio.ModeOut)
	}
	key.Setup(gpio.ModeIn | gpio.PullUp)
}

func main() {
	n := 0
	for {
		dir := 1
		if key.Load() != 0 {
			dir = -1
		}
		n += dir
		switch n {
		case -1:
			n = len(leds) - 1
		case len(leds):
			n = 0
		}
		led := leds[n]
		led.Set()
		delay.Millisec(20)
		led.Clear()
		delay.Millisec(500)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
