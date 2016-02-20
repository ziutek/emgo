package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"
)

type LED struct {
	Port gpio.Port
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

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	gpio.B.EnableClock(true)
	gpio.D.EnableClock(true)

	for _, led := range leds {
		led.Init()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,
}

func wait(ms int) {
	delay.Millisec(ms)
	//delay.Loop(ms * 1e4)
}

func blink(led *LED, d, n int) {
	for n > 0 {
		led.On()
		wait(d)
		led.Off()
		wait(d)
		n--
	}
}

func main() {
	for {
		go blink(&leds[0], 100, 10)
		go blink(&leds[1], 230, 10)
		go blink(&leds[2], 350, 10)
		blink(&leds[3], 500, 10)
		wait(250)
		// BUG: In real application you schould ensure that all gorutines
		// finished before next loop. In this case Blue LED blinks longest
		// so this example works.
	}
}
