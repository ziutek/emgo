package main

import (
	"delay"
	"runtime/noos"
	"strconv"

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

	periph.APB2ClockEnable(periph.SysCfg)
	periph.APB2Reset(periph.SysCfg)
	periph.AHB1ClockEnable(periph.GPIOC)
	periph.AHB1Reset(periph.GPIOC)

	gpio.C.SetMode(1, gpio.In)
	exti.L1.Connect(gpio.C)
	exti.L1.RiseTrigEnable()
	exti.L1.FallTrigEnable()
	exti.L1.IntEnable()
	irqs.Ext1.UseHandler(pulse)
	irqs.Ext1.Enable()

	periph.APB2ClockDisable(periph.SysCfg)
}

func blink(led int) {
	leds.SetBit(led)
	delay.Loop(1e6)
	leds.ClearBit(led)
}

type Bit byte

type Pulse struct {
	Stamp uint64
	Bit   Bit
}

const (
	Zero Bit = 0
	One  Bit = 1
	Sync Bit = 2
	Err  Bit = 3

	ok = Zero
)

var (
	last uint64
	c    = make(chan Pulse, 3)
)

func rising(t uint64) Bit {
	dt64 := t - last
	last = t
	if dt64 > 2050e6 {
		return Err
	}
	dt := uint(dt64)
	switch {
	case dt > 1950e6:
		return Sync
	case dt > 1050e6:
		return Err
	case dt > 950e6:
		return ok
	}
	return Err
}

func falling(t uint64) Bit {
	blen64 := (t - last)
	if blen64 > 250e6 {
		return Err
	}
	blen := uint(blen64)
	switch {
	case blen > 140e6:
		return One
	case blen > 130e6:
		return Err
	case blen > 40e6:
		return Zero
	}
	return Err
}

func pulse() {
	t := noos.Uptime()
	exti.L1.ClearPending()
	if gpio.C.Bit(1) {
		if status := rising(t); status != ok {
			select {
			case c <- Pulse{last, status}:
			default:
			}
		}
		return
	}
	select {
	case c <- Pulse{last, falling(t)}:
	default:
	}
}

func main() {
	for {
		p := <-c
		strconv.WriteUint64(con, p.Stamp, 10)
		switch p.Bit {
		case Zero:
			con.WriteString(" 0\n")
			blink(Blue)
		case One:
			con.WriteString(" 1\n")
			blink(Green)
		case Sync:
			con.WriteString(" sync\n")
			blink(Orange)
		default: // Err
			con.WriteString(" error\n")
			blink(Red)
		}
	}
}
