package main

import (
	"delay"
	"strconv"
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

func blink(led int) {
	leds.SetBit(led)
	delay.Loop(1e6)
	leds.ClearBit(led)
}

var d = dcf77.NewDecoder()

func edgeISR() {
	t := time.Now()
	exti.L1.ClearPending()
	d.Edge(t, gpio.C.Bit(1))
}

func main() {
	for {
		p := d.Pulse()
		strconv.WriteInt64(con, time.Now().UnixNano(), 10)
		con.WriteByte(' ')
		strconv.WriteInt64(con, p.Stamp.UnixNano(), 10)
		con.WriteByte(' ')
		if p.Err != nil {
			blink(Red)
			con.WriteString(p.Err.Error())
		} else {
			blink(Green)
			p.Time.WriteText(con)
		}
		con.WriteByte('\n')
	}
}
