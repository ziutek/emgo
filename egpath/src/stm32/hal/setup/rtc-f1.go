// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl

package setup

import (
	"stm32/hal/raw/pwr"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/rtc"
	//"stm32/hal/raw/bkp"
)

func rtcWaitForSync() {
	rtc.RTC.RSF().Clear()
	for rtc.RTC.RSF().Load() == 0 {
	}
}

func rtcWaitForWrite() {
	for rtc.RTC.RTOFF().Load() == 0 {
	}
}

// When 32768 Hz oscilator is used and rtcPreLog2 == 5 then:
// - rtos.Uptime resolution is 1/32768 s,
// - rtos.SleepUntil resoultion is 1<<5/32768 s = 1/1024 s ≈ 1 ms,
// - the longest down time (RTC on battery) is 1<<32/1024 s ≈ 1165 h ≈ 48 days.
const (
	rtcPreLog2 = 5
	rtcDivMax  = 1<<rtcPreLog2 - 1
)


func useRTC(freq uint) {
	RTC := rtc.RTC
	RCC := rcc.RCC
	PWR := pwr.PWR
	//BKP := bkp.BKP

	const (
		mask = rcc.LSEON | rcc.RTCSEL | rcc.RTCEN
		cfg  = rcc.LSEON | rcc.RTCSEL_LSE | rcc.RTCEN
	)

	if RCC.BDCR.Bits(mask) == cfg {
		// RTC configured before.
		rtcWaitForSync()
		return
	}

	RCC.APB1ENR.SetBits(rcc.PWREN | rcc.BKPEN)
	_ = RCC.APB1ENR.Load()
	PWR.DBP().Set()
	RCC.BDRST().Set()
	RCC.BDRST().Clear()
	RCC.LSEON().Set()
	for RCC.LSERDY().Load() == 0 {
	}
	RCC.BDCR.StoreBits(mask, cfg)

	rtcWaitForSync()
	rtcWaitForWrite()
	RTC.CNF().Set() // Begin RTC configuration
	RTC.PRLH.Store(0)
	RTC.PRLL.Store(rtcDivMax)
	RTC.CNTH.Store(0)
	RTC.CNTL.Store(0)
	RTC.CNF().Clear() // Copy from APB to BKP domain.
	rtcWaitForWrite()

	RCC.APB1ENR.ClearBits(rcc.PWREN | rcc.BKPEN)
}

// rtcNow returns number of ticks of RTC input clock. Returned value is calculated
// using 32-bit counter value and lowest 16 bits of prescaler divider, so rtcNow
// works only fo prescaler <= 65536.
func RtcNow() int64 {
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
	cnt := uint(ch)<<16 | uint(cl)
	frac := rtcDivMax - (uint(dl)+1)&rtcDivMax
	return int64(cnt)<<rtcPreLog2 | int64(frac)
}
