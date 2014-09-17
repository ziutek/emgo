package main

import (
	"delay"
	"fmt"
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
	irqs.Ext1.UseHandler(edgeISR)
	irqs.Ext1.Enable()

	periph.APB2ClockDisable(periph.SysCfg)
}

func blink(led, dly int) {
	leds.SetBit(led)
	if dly < 0 {
		delay.Loop(-dly * 1e4)
	} else {
		delay.Millisec(dly)
	}
	leds.ClearBit(led)
}

var d = dcf77.NewDecoder()

func edgeISR() {
	t := time.Now()
	exti.L1.ClearPending()
	blink(Blue, -50)
	d.Edge(t, gpio.C.Bit(1))
}

func main() {
	for {
		p := d.Pulse()
		now := fmt.Int64(time.Now().UnixNano())
		stamp := fmt.Int64(p.Stamp.UnixNano())
		fmt.Fprint(con, now, fmt.S, stamp, fmt.S)
		if p.Err != nil {
			blink(Red, 50)
			fmt.Fprint(con, fmt.Err(p.Err), fmt.N)
		} else {
			blink(Green, 50)
			fmt.Fprint(con, p.Time, fmt.N)
		}
	}
}
