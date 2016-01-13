package main

import (
	"delay"
	"fmt"

	"stm32/hal/raw/pwr"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/rtc"
	"stm32/hal/setup"
)

func init() {
	setup.Performance(8, 72/8, false)
}

func rtcWaitForSync() {
	rtc.RTC.RSF().Clear()
	for rtc.RTC.RSF().Load() == 0 {
	}
}

func rtcWaitForWrite() {
	for rtc.RTC.RTOFF().Load() == 0 {
	}
}

const (
	rtcPreLog2 = 5
	rtcDivMax  = 1<<rtcPreLog2 - 1
)

// rtcNow returns number of ticks of RTC input clock. Returned value is calculated
// using 32-bit counter value and lowest 16 bits of prescaler divider, so it can
// be used only if prescaler <= 65536.
func rtcNow() int64 {
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

func main() {
	delay.Millisec(100)
	fmt.Println("Start")

	RTC := rtc.RTC
	RCC := rcc.RCC
	PWR := pwr.PWR
	const (
		mask = rcc.LSEON | rcc.RTCSEL | rcc.RTCEN
		cfg  = rcc.LSEON | rcc.RTCSEL_LSE | rcc.RTCEN
	)
	if RCC.BDCR.Bits(mask) != cfg {
		fmt.Println("RTC not initialized. Initializing...")
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

		PWR.DBP().Clear()
		RCC.APB1ENR.ClearBits(rcc.PWREN | rcc.BKPEN)
		fmt.Println("Done.")
	}

	buf := make([]int64, 256)
	var prev int64
	for i := 0; i < len(buf); {
		now := rtcNow()
		if now != prev {
			buf[i] = now
			prev = now
			i++
		}
	}
	for i, v := range buf {
		if i == 0 {
			fmt.Printf("%d: %d\n", i, v)
		} else {
			fmt.Printf("%3d: %d %d\n", i, v, v-buf[i-1])
		}
		delay.Millisec(5)
	}
}
