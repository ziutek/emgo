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

type data struct {
	i           int
	min, hour   int
	mday, wday  int
	month, year int
	ones        int
	summer      bool
}

func (d *data) reset() {
	*d = data{i: -1}
}

func (d *data) zero() {
	*d = data{}
}

func main() {
	var d data
	d.reset()
	for {
		p := <-c

		switch {
		case d.i == -1 || p.Bit > One:
			// break
		case d.i <= 16:
			// ...
		case d.i == 17:
			d.summer = (p.Bit == One)
		case d.i == 18:
			if d.summer && p.Bit == One || !d.summer && p.Bit == Zero {
				con.WriteString("Err: bad summer time bits\n")
				d.reset()
			}
		case d.i == 19:
			// Leap second announcement
		case d.i == 20:
			if p.Bit != 1 {
				con.WriteString("Err: start bit != 1\n")
				d.reset()
			}
		case d.i <= 27:
			if p.Bit == One {
				d.ones++
				d.min += 1 << uint(d.i-21)
			}
		case d.i == 28:
			if p.Bit == One {
				d.ones++
			}
			if d.ones&1 != 0 {
				con.WriteString("Err: minute bits parity != even\n")
				d.reset()
			}
		case d.i <= 34:
			if p.Bit == One {
				d.ones++
				d.hour += 1 << uint(d.i-29)
			}
		case d.i == 35:
			if p.Bit == One {
				d.ones++
			}
			if d.ones&1 != 0 {
				con.WriteString("Err: hour bits parity != even\n")
				d.reset()
			}
		case d.i <= 41:
			if p.Bit == One {
				d.ones++
				d.mday += 1 << uint(d.i-36)
			}
		case d.i <= 44:
			if p.Bit == One {
				d.ones++
				d.wday += 1 << uint(d.i-42)
			}
		case d.i <= 49:
			if p.Bit == One {
				d.ones++
				d.month += 1 << uint(d.i-45)
			}
		case d.i <= 57:
			if p.Bit == One {
				d.ones++
				d.year += 1 << uint(d.i-50)
			}
		case d.i == 58:
			if p.Bit == One {
				d.ones++
			}
			if d.ones&1 != 0 {
				con.WriteString("Err: date bits parity != even\n")
				d.reset()
			}
		}

		strconv.WriteUint64(con, p.Stamp, 10)
		con.WriteByte(' ')
		strconv.WriteInt(con, int32(d.i), 10)
		if d.i >= 0 {
			d.i++
		}
		con.WriteByte(' ')
		switch p.Bit {
		case Zero:
			con.WriteByte('0')
			blink(Blue)
		case One:
			con.WriteByte('1')
			blink(Green)
		case Sync:
			con.WriteString("sync\n")
			blink(Orange)
			con.WriteString("Time: ")
			strconv.WriteInt(con, int32(2<<12+d.year), 16)
			con.WriteByte('-')
			strconv.WriteInt(con, int32(d.month), 16)
			con.WriteByte('-')
			strconv.WriteInt(con, int32(d.mday), 16)
			con.WriteByte(' ')
			strconv.WriteInt(con, int32(d.hour), 16)
			con.WriteByte(':')
			strconv.WriteInt(con, int32(d.min), 16)
			con.WriteByte(' ')
			if d.summer {
				con.WriteString("CEST")
			} else {
				con.WriteString("CET")
			}
			d.zero()
		default: // Err
			con.WriteString("error")
			blink(Red)
			d.reset()
		}
		con.WriteByte('\n')
	}
}
