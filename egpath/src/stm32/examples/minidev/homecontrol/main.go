package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"

	"stm32/hal/raw/afio"
	"stm32/hal/raw/rcc"
)

var (
	led    gpio.Pin
	relays [4]gpio.Pin
	ssr    gpio.Pin
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	// Allocate pins.

	gpio.B.EnableClock(false)
	relays[3] = gpio.B.Pin(4)
	relays[2] = gpio.B.Pin(5)
	relays[1] = gpio.B.Pin(6)
	relays[0] = gpio.B.Pin(7)
	ssr = gpio.B.Pin(8)

	gpio.C.EnableClock(false)
	led = gpio.C.Pin(13)

	// Configure pins.

	// Release JTDI and NJTRST (PA15 and PB4) to use as GPIO pins.
	rcc.RCC.AFIOEN().Set()
	afio.AFIO.SWJ_CFG().Store(afio.SWJ_CFG_JTAGDISABLE)
	rcc.RCC.AFIOEN().Clear()

	cfg := &gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low}
	led.Setup(cfg)
	for _, pin := range relays {
		pin.Set()
		pin.Setup(cfg)
	}
	cfg = &gpio.Config{Mode: gpio.Out, Driver: gpio.PushPull, Speed: gpio.Low}
	ssr.Setup(cfg)
}

func main() {
	for _, relay := range relays {
		led.Clear()
		relay.Clear()
		delay.Millisec(50)
		led.Set()
		delay.Millisec(950)
	}
	ssr.Set()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,
}
