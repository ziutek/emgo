package main

import (
	"delay"
	"rtos"

	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var (
	led, btn gpio.Pin
	ch       = make(chan struct{}, 3)
)

func init() {
	system.Setup80(0, 0)
	systick.Setup(2e6)

	gpio.A.EnableClock(false)
	led = gpio.A.Pin(5)
	gpio.C.EnableClock(true)
	btn = gpio.C.Pin(13)

	led.Setup(&gpio.Config{Mode: gpio.Out, Speed: gpio.Low})
	btn.Setup(&gpio.Config{Mode: gpio.In})
	li := exti.Lines(btn.Mask())
	li.Connect(btn.Port())
	li.EnableFallTrig()
	li.EnableIRQ()

	rtos.IRQ(irq.EXTI15_10).Enable()
}

func blink(ms int) {
	led.Set()
	delay.Millisec(ms)
	led.Clear()
}

func main() {
	for {
		select {
		case <-ch:
			blink(50)
		case <-rtos.At(rtos.Nanosec() + 2e9):
			blink(250)
		}
	}
}

func btnISR() {
	exti.Lines(btn.Mask()).ClearPending()
	select {
	case ch <- struct{}{}:
	default:
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI15_10: btnISR,
}
