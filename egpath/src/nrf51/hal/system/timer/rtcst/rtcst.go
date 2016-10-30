// Package rtcst implements tickless system timer using real time counter.
package rtcst

import (
	"rtos"
	"sync/atomic"
	"syscall"

	"nrf51/hal/rtc"
	"nrf51/hal/te"
)

type globals struct {
	st  *rtc.Periph
	cce *te.Event
	cnt uint32
	mul uint32
}

var g globals

// Setup setups st as system timer using ccn compare channel number. Usually
// rtc.RTC0 and channel number 1 is used. Setup accepts st at its reset state or
// configured before for other purposes. For its needs it uses only rtc.OVRFLW
// and rtc.COMPARE0+ccn events and can work with any prescaler set before
// (avoid prescaler values that cause tick period > 1 ms). Setup starts st by
// triggering rtc.START task but accepts RTC started before. It setups st.IRQ()
// priority in NVIC and enables IRQ handling.
func Setup(st *rtc.Periph, ccn int) {
	if uint(ccn) > 3 {
		panic("rtcst: bad ccn")
	}
	g.st = st
	g.cce = st.Event(rtc.COMPARE0 + rtc.Event(ccn))
	g.cce.DisablePPI()
	g.cce.DisableIRQ()
	g.mul = (st.PRESCALER() + 1) * (1e9 >> 9)
	ove := st.Event(rtc.OVRFLW)
	ove.DisablePPI()
	ove.EnableIRQ()
	st.Task(rtc.START).Trigger()
	irq := rtos.IRQ(st.IRQ())
	irq.SetPrio(rtos.IRQPrioLowest + rtos.IRQPrioStep*rtos.IRQPrioNum*3/4)
	irq.Enable()
	syscall.SetSysTimer(Nanosec, nil)
}

func ISR() {
	ove := g.st.Event(rtc.OVRFLW)
	cce := g.cce
	switch {
	case ove.IsSet():
		ove.Clear()
		g.cnt++
	case cce.IsSet():
		cce.Clear()
	}
}

func ticktons(ch, cl uint32) int64 {
	// Exact and efficient calculation of: (ch<<24+cl)*1e9*(prescaler+1)/32768.
	mul := int64(g.mul)
	h := int64(ch) * mul << (24 - (15 - 9))
	l := int64(cl) * mul >> (15 - 9)
	return h + l
}

func Nanosec() int64 {
	softcnt := atomic.LoadUint32(&g.cnt)
	for {
		rtccnt := g.st.COUNTER()
		softcnt1 := atomic.LoadUint32(&g.cnt)
		if softcnt1 == softcnt {
			return ticktons(softcnt, rtccnt)
		}
		softcnt = softcnt1
	}
	// BUG?: Can Nanosec read COUNTER, g.cnt after overflow but before ISR?
}
