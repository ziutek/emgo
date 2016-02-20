// This example doesnt work! There is a BUG somewhere in exti or ...
package main

import (
	"delay"
	"rtos"

	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"
)

var leds, keys gpio.Port

const (
	LED1 = gpio.Pin7
	LED2 = gpio.Pin6
	LED3 = gpio.Pin5

	Key3 = gpio.Pin10
)

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	gpio.B.EnableClock(false)
	leds = gpio.B
	gpio.C.EnableClock(true)
	keys = gpio.C

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(LED1|LED2|LED3, &cfg)

	cfg = gpio.Config{Mode: gpio.In, Pull: gpio.PullUp}
	keys.Setup(Key3, &cfg)
	line := exti.Lines(Key3)
	line.Connect(keys)
	line.EnableFallTrig()
	line.EnableInt()

	rtos.IRQ(irq.EXTI15_10).Enable()
}

func keyISR() {
	exti.Lines(Key3).ClearPending()
	leds.SetPins(LED3)
	delay.Loop(1e5)
	leds.ClearPins(LED3)
}

func main() {
	for {
		leds.SetPins(LED1)
		delay.Millisec(1000)
		leds.ClearPins(LED1)
		delay.Millisec(1000)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm:  rtc.ISR,
	irq.EXTI15_10: keyISR,
}
