package main

import (
	"delay"
	"fmt"
	"rtos"
	"time"

	"dcf77"

	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/setup"
)

const dcfpin = gpio.Pin1

var dcfport *gpio.Port

func init() {
	setup.Performance168(8)
	setup.UseSysTick()

	gpio.C.EnableClock(true)
	dcfport = gpio.C
	gpio.D.EnableClock(false)
	initLEDs(gpio.D)

	line := exti.Lines(dcfpin)
	line.Connect(dcfport)
	line.EnableRiseTrig()
	line.EnableFallTrig()
	line.EnableInt()

	rtos.IRQ(irq.EXTI1).Enable()
}

var d = dcf77.NewDecoder()

func edgeISR() {
	t := time.Now()
	exti.Lines(dcfpin).ClearPending()
	blink(Blue, -100)
	d.Edge(t, dcfport.Pins(dcfpin) != 0)
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI1: edgeISR,
}

func main() {
	delay.Millisec(100)
	for {
		p := d.Pulse()
		now := time.Now().UnixNano()
		if p.Err() != nil {
			fmt.Printf("now=%d %v\n", now, p.Err())
			blink(Red, 25)
		} else {
			stamp := p.Stamp.UnixNano()
			fmt.Printf("now=%d stamp=%d dcf=%s\n", now, stamp, p.Date)
			blink(Green, 25)
		}
	}
}
