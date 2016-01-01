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

const dcfpin = 1

var dcfport = gpio.C

func init() {
	setup.Performance168(8)

	initLEDs()

	// Initialize DCF77 input pin.

	dcfport.EnableClock(true)
	dcfport.SetMode(dcfpin, gpio.In)

	line := exti.Line(dcfpin)
	line.Connect(dcfport)
	line.EnableRiseTrig()
	line.EnableFallTrig()
	line.EnableInt()

	rtos.IRQ(irq.EXTI1).Enable()
}

var d = dcf77.NewDecoder()

func edgeISR() {
	t := time.Now()
	exti.Line(dcfpin).ClearPending()
	blink(Blue, -100)
	d.Edge(t, dcfport.PinIn(dcfpin) != 0)
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI1: edgeISR,
}

func main() {
	delay.Millisec(100)
	for {
		fmt.Println("Pulse wait...")
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
