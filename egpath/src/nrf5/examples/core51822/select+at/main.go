// This example shows how to use GPIOTE peripheral to detect changes at specific
// GPIO pin and handle them using interrupt. Additionaly it shows how to use
// asynchronous (buffered) channel for communication between ISR and gorutine
// (task/thread) and how to use rtos.At function to provide deadline/timeout
// for such communication.
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

const key = gpiote.Chan(0)

var (
	leds = make([]gpio.Pin, 5)
	ch   = make(chan struct{}, 1)
)

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// Allocate pins (always do it in one place to avoid conflicts).
	p0 := gpio.P0
	keyPin := p0.Pin(16) // KEY1
	for i := range leds {
		leds[i] = p0.Pin(18 + i)
	}

	// Configure pins.
	for _, led := range leds {
		led.Setup(gpio.ModeOut)
	}
	keyPin.Setup(gpio.ModeIn | gpio.PullUp)

	// Configure GPIOTE
	key.Setup(keyPin, gpiote.ModeEvent|gpiote.PolarityHiToLo)
	key.IN().Event().EnableIRQ()

	// Enable IRQs in NVIC.
	rtos.IRQ(irq.GPIOTE).Enable()
}

func blinkNextLED() {
	n := len(leds) - 1
	leds[n].Set()
	delay.Millisec(50)
	leds[n].Clear()
	if n > 0 {
		leds = leds[:n]
	} else {
		leds = leds[:cap(leds)]
	}
}

func main() {
	for {
		select {
		case <-ch:
		case <-rtos.At(rtos.Nanosec() + 2e9):
		}
		blinkNextLED()
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
