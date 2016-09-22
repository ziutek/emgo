package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"
)

var leds *gpio.Port

const Blue = gpio.Pin13

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	gpio.C.EnableClock(true)
	leds = gpio.C

	leds.Setup(Blue, gpio.Config{Mode: gpio.Out, Speed: gpio.Low})
}

func main() {
	for {
		for _, d := range []int{100, 100, 400, 400, 400, 400, 100, 100, 100} {
			leds.ClearPins(Blue)
			delay.Millisec(20)
			leds.SetPins(Blue)
			delay.Millisec(d)
		}
		delay.Millisec(1000)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,
}
