// Package rtcst implements a tickless system timer using the real time counter.
package rtcst

import (
	"math"
	"rtos"
	"sync/atomic"
	"sync/fence"
	"syscall"

	"nrf5/hal/rtc"
	"nrf5/hal/te"
)

type globals struct {
	wakens  int64
	st      *rtc.Periph
	softcnt uint32
	scale   uint32
	ccn     byte
	alarm   bool
}

var g globals

func cce() *te.Event {
	return g.st.Event(rtc.COMPARE0 + rtc.Event(g.ccn))
}

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
	g.ccn = byte(ccn)
	cce := cce()
	cce.DisablePPI()
	cce.DisableIRQ()
	g.scale = (st.PRESCALER() + 1) * (1e9 >> 9)
	ove := st.Event(rtc.OVRFLW)
	ove.DisablePPI()
	ove.EnableIRQ()
	irq := rtos.IRQ(st.IRQ())
	// Priority of rtc.OVRFLW IRQ must be higher than SVCall proprity, to ensure
	// that any user of rtos.Nanosec observes both COUNTER reset and softcnt
	// increment as one atomic operation.
	irq.SetPrio(rtos.IRQPrioLowest + rtos.IRQPrioStep*rtos.IRQPrioNum*3/4)
	irq.Enable()
	st.Task(rtc.START).Trigger()
	syscall.SetSysTimer(nanosec, setWakeup)
}

func ISR() {
	ove := g.st.Event(rtc.OVRFLW)
	cce := cce()
	if ove.IsSet() {
		ove.Clear()
		atomic.StoreUint32(&g.softcnt, g.softcnt+1)
	}
	if cce.IsSet() {
		cce.DisableIRQ()
		g.wakens = 0
		if g.alarm {
			g.alarm = false
			syscall.Alarm.Send()
		} else {
			syscall.SchedNext()
		}
	}
}

func cntstons(ch, cl uint32) int64 {
	// Exact and efficient calculation of: (ch<<24+cl)*1e9*(prescaler+1)/32768.
	scale := int64(g.scale)
	h := int64(ch) * scale << (24 - (15 - 9))
	l := int64(cl) * scale >> (15 - 9)
	return h + l
}

func counters() (ch, cl uint32) {
	ch = atomic.LoadUint32(&g.softcnt)
	for {
		fence.R() // Ensure IO load after load(g.softcnt).
		cl = g.st.COUNTER()
		fence.R() // Ensure load(g.softcnt) after IO load.
		ch1 := atomic.LoadUint32(&g.softcnt)
		if ch1 == ch {
			return
		}
		ch = ch1
	}
}

func ticks(ch, cl uint32) int64 {
	return int64(ch)<<24 | int64(cl)
}

// nanosec: see syscall.SetSysTimer.
func nanosec() int64 {
	return cntstons(counters())
}

func nstotick(ns int64) int64 {
	return int64(math.Muldiv(uint64(ns), uint64(1<<(15-9)), uint64(g.scale)))
}

// setWakeup: see syscall.SetSysTimer.
func setWakeup(ns int64, alarm bool) {
	if g.wakens == ns && g.alarm == alarm {
		return
	}
	g.wakens = ns
	g.alarm = alarm
	wkup := nstotick(ns) + 1 // +1 to don't wakeup to early because of rounding.
	cce := cce()
	cce.Clear()
	ch, cl := counters()
	now := ticks(ch, cl)
	sleep := wkup - now
	switch {
	case sleep > 0xffffff:
		g.st.SetCC(int(g.ccn), cl)
	case sleep > 0:
		g.st.SetCC(int(g.ccn), uint32(wkup)&0xffffff)
		now = ticks(counters())
	}
	if now < wkup {
		cce.EnableIRQ()
		return
	}
	// wkup in the past or there is a chance that CC was set to late.
	g.wakens = 0
	if g.alarm {
		g.alarm = false
		syscall.Alarm.Send()
	} else {
		syscall.SchedNext()
	}
}
