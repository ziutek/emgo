package main

import (
	"delay"
	"fmt"
	"time"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"
)

var leds *gpio.Port

const (
	LED1 = gpio.Pin7
	LED2 = gpio.Pin6
)

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	gpio.B.EnableClock(true)
	leds = gpio.B

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(LED1|LED2, cfg)
}


func printDate(led gpio.Pins, dly int) {
	for {
		leds.SetPins(led)
		delay.Millisec(dly)
		leds.ClearPins(led)
		delay.Millisec(dly)
		t := time.Now()
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		ns := t.Nanosecond()
		fmt.Printf(
			"%04d-%02d-%02d %02d:%02d:%02d.%09d\n",
			y, mo, d, h, mi, s, ns,
		)
	}
}

func main() {
	if ok, set := rtc.Status(); ok && !set {
		rtc.SetTime(time.Date(2016, 1, 24, 22, 58, 30, 0, time.UTC))
	}
	printDate(LED2, 500)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,
}
