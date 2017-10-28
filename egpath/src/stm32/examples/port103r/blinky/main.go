package main

import (
	"delay"

	"stm32/hal/gpio"
	//"stm32/hal/irq"
	"stm32/hal/system"
	//"stm32/hal/system/timer/rtcst"
	"stm32/hal/system/timer/systick"
)

var leds *gpio.Port

const (
	LED1 = gpio.Pin7
	LED2 = gpio.Pin6
	LED3 = gpio.Pin5
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	systick.Setup(2e6)
	//rtcst.Setup(32768)

	gpio.B.EnableClock(true)
	leds = gpio.B

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(LED1|LED2|LED3, &cfg)
}

func blink(led gpio.Pins, dly int) {
	for {
		leds.SetPins(led)
		delay.Millisec(dly)
		leds.ClearPins(led)
		delay.Millisec(dly)
	}
}

func main() {
	//go blink(LED1, 500)
	//blink(LED2, 1000)
	
	for {
		leds.SetPins(LED3)
		delay.Millisec(2000)
		leds.ClearPins(LED3)
		delay.Millisec(2000)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	//irq.RTCAlarm: rtcst.ISR,
}
