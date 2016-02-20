package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"
)

type LED struct {
	Port *gpio.Port
	Pin  gpio.Pins
}

func (led LED) Init() {
	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	led.Port.Setup(led.Pin, cfg)
}

func (led LED) On() {
	led.Port.SetPins(led.Pin)
}

func (led LED) Off() {
	led.Port.ClearPins(led.Pin)
}

var leds = []LED{
	{gpio.B, gpio.Pin7},
	{gpio.B, gpio.Pin6},
	{gpio.B, gpio.Pin5},
	{gpio.D, gpio.Pin2},
}

var keys *gpio.Port

const key3 = gpio.Pin10

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	gpio.B.EnableClock(true)
	gpio.C.EnableClock(true)
	gpio.D.EnableClock(true)

	for _, led := range leds {
		led.Init()
	}
	keys = gpio.C
	keys.Setup(key3, &gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,
}

func main() {
	i := 0
	d := -1
	n := len(leds) - 1
	for {
		leds[i].On()
		delay.Millisec(50)
		leds[i].Off()
		if i == 0 || i == n || keys.Pins(key3) == 0 {
			d = -d
		}
		i += d
	}
}
