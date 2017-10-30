package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
)

var (
	led1 = gpio.B.Pin(7)
	led2 = gpio.B.Pin(6)
	led3 = gpio.B.Pin(5)
	led4 = gpio.D.Pin(2)
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	gpio.B.EnableClock(false)
	gpio.D.EnableClock(false)

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	led1.Setup(&cfg)
	led2.Setup(&cfg)
	led3.Setup(&cfg)
	led4.Setup(&cfg)
}

func main() {
	for {
		led4.Clear()
		led1.Set()
		delay.Millisec(500)
		led1.Clear()
		led2.Set()
		delay.Millisec(500)
		led2.Clear()
		led3.Set()
		delay.Millisec(500)
		led3.Clear()
		led4.Set()
		delay.Millisec(500)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,
}
