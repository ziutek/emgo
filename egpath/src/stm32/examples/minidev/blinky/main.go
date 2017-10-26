package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
)

var leds *gpio.Port

const Blue = gpio.Pin13

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	gpio.C.EnableClock(true)
	leds = gpio.C

	cfg := gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low}
	leds.Setup(Blue, &cfg)
}

func main() {
	for {
		leds.ClearPins(Blue)
		delay.Millisec(50)
		leds.SetPins(Blue)
		delay.Millisec(2950)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,
}
