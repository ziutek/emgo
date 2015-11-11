package main

import (
	"delay"
	"fmt"
	"rtos"
	"time"

	"dcf77"
	"stm32/f4/exti"
	"stm32/f4/gpio"
	"stm32/f4/irqs"
	"stm32/f4/periph"
	"stm32/f4/setup"
)

func init() {
	setup.Performance168(8)

	initLEDs()
	initConsole()

	// Initialize DCF77 input pin.

	periph.APB2ClockEnable(periph.SysCfg)
	periph.APB2Reset(periph.SysCfg)
	periph.AHB1ClockEnable(periph.GPIOC)
	periph.AHB1Reset(periph.GPIOC)

	gpio.C.SetMode(1, gpio.In)
	exti.L1.Connect(gpio.C)
	exti.L1.RiseTrigEnable()
	exti.L1.FallTrigEnable()
	exti.L1.IntEnable()
	rtos.IRQ(irqs.Ext1).UseHandler(edgeISR)
	rtos.IRQ(irqs.Ext1).Enable()

	periph.APB2ClockDisable(periph.SysCfg)
}

func blink(led int, dly int) {
	leds.SetPin(led)
	if dly < 0 {
		delay.Loop(-dly * 1e3)
	} else {
		delay.Millisec(dly)
	}
	leds.ClearPin(led)
}

var d = dcf77.NewDecoder()

func edgeISR() {
	t := time.Now()
	exti.L1.ClearPending()
	blink(Blue, -100)
	d.Edge(t, gpio.C.InPin(1) != 0)
}

func main() {
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
