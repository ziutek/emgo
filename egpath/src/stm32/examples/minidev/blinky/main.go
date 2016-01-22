package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/setup"
)

var leds *gpio.Port

const Blue = gpio.Pin13

func init() {
	setup.Performance(8, 72/8, false)
	setup.UseRTC(32768)

	gpio.C.EnableClock(true)
	leds = gpio.C

	leds.Setup(Blue, &gpio.Config{Mode: gpio.Out, Speed: gpio.Low})
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: setup.RTCISR,
}

func wait(ms int) {
	//delay.Loop(1e7)
	delay.Millisec(ms)
}

func main() {
	for {
		for _, d := range []int{100, 100, 400, 400, 400, 400, 100, 100, 100} {
			leds.ClearPins(Blue)
			wait(20)
			leds.SetPins(Blue)
			wait(d)
		}
		wait(1000)
	}
}
