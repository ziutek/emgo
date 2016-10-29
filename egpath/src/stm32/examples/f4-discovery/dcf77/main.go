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
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

const dcfpin = gpio.Pin1

var dcfport *gpio.Port

func init() {
	system.Setup168(8)
	systick.Setup()

	gpio.C.EnableClock(true)
	dcfport = gpio.C
	gpio.D.EnableClock(false)
	initLEDs(gpio.D)

	line := exti.Lines(dcfpin)
	line.Connect(dcfport)
	line.EnableRisiTrig()
	line.EnableFallTrig()
	line.EnableIRQ()

	prio := (rtos.IRQPrioLowest + rtos.SyscallPrio) / 2
	rtos.IRQ(irq.EXTI1).SetPrio(prio) // To allow use time.Now().
	rtos.IRQ(irq.EXTI1).Enable()
}

var d = dcf77.NewDecoder()

func edgeISR() {
	t := time.Now()
	exti.Lines(dcfpin).ClearPending()
	blink(Blue, -1e4)
	d.Edge(t, dcfport.Pins(dcfpin) != 0)
}

func main() {
	delay.Millisec(100)
	fmt.Printf("\nDCF77 RECEIVER\n\n")
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

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI1: edgeISR,
}
