// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl

package setup

import (
	"math"
	"mmio"
	"rtos"
	"sync/atomic"
	"syscall"

	"arch/cortexm/scb"

	"stm32/hal/exti"
	"stm32/hal/irq"

	"stm32/hal/raw/bkp"
	"stm32/hal/raw/pwr"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/rtc"
)

// When 32768 Hz oscilator is used and rtcPreLog2 == 5 then:
// - rtos.Nanosec resolution will be 1/32768 s,
// - rtos.SleepUntil resoultion will be 1<<5/32768 s = 1/1024 s ≈ 1 ms,
// - the longest down time (RTC on battery) will be 1<<32/1024 s ≈ 48 days.
const (
	rtcPreLog2 = 5
	rtcPre     = 1 << rtcPreLog2

	rtcMaxSleepCnt  = 1 << 22
	rtcMaxSleepTick = rtcMaxSleepCnt << rtcPreLog2
)

const (
	rtcDirty = iota
	rtcNotSet
	rtcOK
)

func ticktons(tick int64) int64 {
	return int64(math.Muldiv(uint64(tick), 1e9, uint64(rtcFreqHz)))
}

func nstotick(ns int64) int64 {
	return int64(math.Muldiv(uint64(ns), uint64(rtcFreqHz), 1e9))
}

func waitForSync(RTC *rtc.RTC_Periph) {
	RTC.RSF().Clear()
	for RTC.RSF().Load() == 0 {
	}
}

func waitForWrite(RTC *rtc.RTC_Periph) {
	for RTC.RTOFF().Load() == 0 {
	}
}

type twoReg struct {
	high, low *mmio.U16
}

func (tr twoReg) Load() uint32 {
	return uint32(tr.high.Load())<<16 | uint32(tr.low.Load())
}

func (tr twoReg) Store(u uint32) {
	tr.high.Store(uint16(u >> 16))
	tr.low.Store(uint16(u))
}

type rtcBackup struct {
	p *bkp.BKP_Periph
}

func (b rtcBackup) Status() *mmio.U16 {
	return &b.p.DR1.U16
}

func (b rtcBackup) CntExt() *mmio.U16 {
	return &b.p.DR2.U16
}

func (b rtcBackup) LastISR() twoReg {
	return twoReg{&b.p.DR3.U16, &b.p.DR4.U16}
}

var (
	rtcFreqHz uint

	rtcStatus  uint16
	rtcCntExt  uint32
	rtcLastISR uint32
)

func useRTC(freqHz uint) {
	rtcFreqHz = freqHz

	RTC := rtc.RTC
	RCC := rcc.RCC
	PWR := pwr.PWR
	bkp := rtcBackup{bkp.BKP}

	const (
		mask = rcc.LSEON | rcc.RTCSEL | rcc.RTCEN
		cfg  = rcc.LSEON | rcc.RTCSEL_LSE | rcc.RTCEN
	)

	// Enable write access to backup domain.
	RCC.APB1ENR.SetBits(rcc.PWREN | rcc.BKPEN)
	_ = RCC.APB1ENR.Load()
	PWR.DBP().Set()
	RCC.APB1ENR.ClearBits(rcc.PWREN)

	rtcStatus = bkp.Status().Load()

	if RCC.BDCR.Bits(mask) != cfg || rtcStatus == rtcDirty {
		// RTC not initialized or in dirty state.

		// Reset backup domain and configure RTC clock source.
		RCC.BDRST().Set()
		RCC.BDRST().Clear()
		RCC.LSEON().Set()
		for RCC.LSERDY().Load() == 0 {
		}
		RCC.BDCR.StoreBits(mask, cfg)

		// Configure RTC prescaler.
		waitForSync(RTC)
		waitForWrite(RTC)
		RTC.CNF().Set() // Begin PRL configuration
		RTC.PRLL.Store(rtcPre - 1)
		RTC.CNF().Clear() // Copy from APB to BKP domain.

		rtcStatus = rtcNotSet
		bkp.Status().Store(rtcStatus)

		// Wait for complete before setup RTCALR interrupt.
		waitForWrite(RTC)
	} else {
		rtcCntExt = uint32(bkp.CntExt().Load())
		rtcLastISR = bkp.LastISR().Load()
	}
	// Wait for sync. Need in both cases: after reset (synchronise APB domain)
	// or after configuration (avoid reading bad DIVL).
	waitForSync(RTC)

	exti.RTCALR.EnableRiseTrig()
	exti.RTCALR.EnableInt()
	spnum := rtos.IRQPrioStep * rtos.IRQPrioNum
	rtos.IRQ(irq.RTCAlarm).SetPrio(rtos.IRQPrioLowest + spnum*3/4)
	rtos.IRQ(irq.RTCAlarm).Enable()

	// Force RTCISR to early handle possible overflow.
	exti.RTCALR.Trigger()

	syscall.SetSysClock(rtcNanosec, rtcSetWakeup)
}

// rtcVCNT returns value of virtual counter that counts number of ticks of
// RTC input clock. Value of this virtual counter is calculated according to
// the formula:
//
//  VCNT = ((CNTH<<16 + CNTL)<<rtcPreLog2 + frac) & (rtcPre<<32 - 1)
//
// where frac is calculated as follow:
//
//  frac = rtcPre - (DIVL+1)&(rtcPre-1)
//
// Only DIVL is used, so prescaler (rtcPre) can not be greater than 0x10000.
//
// Thanks to this transformation, RTC interrupts are generated at right time.
// See example for Second, Overflow and Alarm(0-1) interrupts in table below:
//
//  CNT      DIV| vcnt>>5 vcnt&0x1f
//  ------------+--------------------
//  ffffffff 04 | ffffffff 1b
//  ffffffff 03 | ffffffff 1c
//  ffffffff 02 | ffffffff 1d
//  ffffffff 01 | ffffffff 1e
//  ffffffff 00 | ffffffff 1f
//  ffffffff 1f | 00000000 00 <- Second, Overflow, Alarm(0xffffffff)
//  00000000 1e | 00000000 01
//  00000000 1d | 00000000 02
//  00000000 1c | 00000000 03
//  00000000 1b | 00000000 04
//  00000000 1a | 00000000 05
//
func rtcVCNT() int64 {
	RTC := rtc.RTC
	var (
		ch rtc.CNTH_Bits
		cl rtc.CNTL_Bits
		dl rtc.DIVL_Bits
	)
	ch = RTC.CNTH.Load()
	for {
		cl = RTC.CNTL.Load()
		for {
			dl = RTC.DIVL.Load()
			cl1 := RTC.CNTL.Load()
			if cl1 == cl {
				break
			}
			cl = cl1
		}
		ch1 := RTC.CNTH.Load()
		if ch1 == ch {
			break
		}
		ch = ch1
	}
	cnt := uint32(ch)<<16 | uint32(cl)
	frac := rtcPre - (uint32(dl)+1)&(rtcPre-1)
	return (int64(cnt)<<rtcPreLog2 + int64(frac)) & (rtcPre<<32 - 1)
}

func RTCISR() {
	exti.RTCALR.ClearPending()

	vcnt32 := uint32(rtcVCNT() >> rtcPreLog2)
	if vcnt32 != rtcLastISR {
		if vcnt32 < rtcLastISR {
			cntExt := rtcCntExt + 1
			bkp := rtcBackup{bkp.BKP}
			bkp.Status().Store(rtcDirty)
			bkp.CntExt().Store(uint16(cntExt))
			bkp.LastISR().Store(vcnt32)
			bkp.Status().Store(rtcStatus)
			atomic.StoreUint32(&rtcCntExt, cntExt)
		} else {
			rtcBackup{bkp.BKP}.LastISR().Store(vcnt32)
		}
		atomic.StoreUint32(&rtcLastISR, vcnt32)
	}

	scb.ICSR_Store(scb.PENDSVSET)
}

func rtcTicks() int64 {
	irq.RTCAlarm.Disable()
	lastISR := atomic.LoadUint32(&rtcLastISR)
	cntExt := atomic.LoadUint32(&rtcCntExt)
	vcnt := rtcVCNT()
	irq.RTCAlarm.Enable()

	if uint32(vcnt>>rtcPreLog2) < lastISR {
		cntExt++
	}
	return int64(cntExt)<<(32+rtcPreLog2) | vcnt
}

// rtcNanosec: see syscall.SetSysClock.
func rtcNanosec() int64 {
	return ticktons(rtcTicks())
}

// rtcSetWakeup: see syscall.SetSysClock.
func rtcSetWakeup(ns int64) {
	wkup := nstotick(ns)
	now := rtcTicks()
	sleep := wkup - now
	switch {
	case sleep > rtcMaxSleepTick:
		wkup = now + rtcMaxSleepTick
	case sleep < 0:
		wkup = now
	}
	wkup = (wkup + rtcPre/2) >> rtcPreLog2
	alr := uint32(wkup) - 1

	RTC := rtc.RTC
	waitForWrite(RTC)
	RTC.CNF().Set()
	RTC.ALRH.Store(rtc.ALRH_Bits(alr >> 16))
	RTC.ALRL.Store(rtc.ALRL_Bits(alr))
	RTC.CNF().Clear()

	if rtcTicks()>>rtcPreLog2 >= wkup {
		// There is a chance that the alarm interrupt was not triggered.
		exti.RTCALR.Trigger()
	}
	/*print64("*sw* now:", now)
	print64(" wkup:", wkup)
	print32(" cnt1:", uint32(now>>rtcPreLog2))
	println32(" alr:", alr)*/
}

/*
const dbg = itm.Port(17)

func print64(s string, i int64) {
	dbg.WriteString(s)
	strconv.WriteInt64(dbg, i, 16, 0)
}

func println64(s string, i int64) {
	dbg.WriteString(s)
	strconv.WriteInt64(dbg, i, 16, 0)
	dbg.WriteString("\r\n")
}

func print32(s string, u uint32) {
	dbg.WriteString(s)
	strconv.WriteUint32(dbg, u, 16, 0)
}

func println32(s string, u uint32) {
	dbg.WriteString(s)
	strconv.WriteUint32(dbg, u, 16, 0)
	dbg.WriteString("\r\n")
}
*/
