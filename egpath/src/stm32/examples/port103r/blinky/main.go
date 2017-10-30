package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
)

var leds = gpio.B

const (
	LED1 = gpio.Pin7
	LED2 = gpio.Pin6
	LED3 = gpio.Pin5
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

	gpio.B.EnableClock(true)
	gpio.D.EnableClock(true)

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	led1.Setup(&cfg)
	led2.Setup(&cfg)
	led3.Setup(&cfg)
	led4.Setup(&cfg)
}

func main() {
	for {
		led4.Set()
		delay.Millisec(2000)
		led4.Clear()
		delay.Millisec(2000)
	}
}

var isrLED = 1

func rtcstISR() {
	led1.Store(isrLED)
	//isrLED ^= 1
	rtcst.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcstISR,
}
