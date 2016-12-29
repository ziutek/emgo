// This example shows how to use channels to divide interrupt handler into two
// parts: fast part - that runs in interrupt context and slow part - that runs
// in thread context.
//
// Fast part only enqueues events/data that receives from source (you) and
// informs the source (using LED3) if its buffer is full. Slow part dequeues
// events/data and handles them. This scheme can be used to implement interrupt
// driven I/O library.
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

var leds, keys *gpio.Port

const (
	LED1 = gpio.Pin7
	LED2 = gpio.Pin6
	LED3 = gpio.Pin5

	Key3 = gpio.Pin10
)

func init() {
	system.Setup(8, 1, 72/8)
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
	line.EnableIRQ()

	rtos.IRQ(irq.EXTI15_10).Enable()
}

var (
	c   = make(chan gpio.Pins, 3)
	led = LED1
)

func keyISR() {
	exti.Lines(Key3).ClearPending()
	select {
	case c <- led:
		if led == LED1 {
			led = LED2
		} else {
			led = LED1
		}
	default:
		// Signal that c is full.
		leds.SetPins(LED3)
		delay.Loop(1e5)
		leds.ClearPins(LED3)
	}
}

func toggle(pins gpio.Pins) {
	leds.SetPins(pins)
	delay.Millisec(500)
	leds.ClearPins(pins)
	delay.Millisec(500)
}

func main() {
	for {
		toggle(<-c)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm:  rtc.ISR,
	irq.EXTI15_10: keyISR,
}
