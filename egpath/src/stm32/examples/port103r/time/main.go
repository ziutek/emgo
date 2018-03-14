package main

import (
	"delay"
	"fmt"
	"rtos"
	"time"
	"time/tz"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
)

var leds *gpio.Port

const (
	LED1 = gpio.Pin7
	LED2 = gpio.Pin6
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	gpio.B.EnableClock(true)
	leds = gpio.B

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(LED1|LED2, &cfg)
}

func main() {
	ok, set := rtcst.Status()
	for !ok {
		fmt.Printf("RTC error\n")
		delay.Millisec(1000)
	}
	if set {
		time.Local = &tz.EuropeWarsaw
	} else {
		t := time.Date(2018, 3, 25, 1, 59, 50, 0, &tz.EuropeWarsaw)
		rtcst.SetTime(t, rtos.Nanosec())
	}
	for {
		leds.SetPins(LED1)
		delay.Millisec(500)
		leds.ClearPins(LED1)
		delay.Millisec(500)
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

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,
}
