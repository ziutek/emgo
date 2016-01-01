package main

import (
	"delay"
	"fmt"

	"stm32/hal/gpio"
	"stm32/hal/setup"

	"stm32/hal/raw/pwr"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/rtc"
)

var LEDport = gpio.A

const Green = 5

func init() {
	setup.Performance96(8)
	LEDport.EnableClock(false)
	LEDport.SetMode(Green, gpio.Out)
}

func wait() {
	delay.Millisec(250)
}

func main() {
	PWR := pwr.PWR
	RCC := rcc.RCC
	RTC := rtc.RTC

	const lse = 1 * rcc.RTCSEL_0
	const bdcrcfg = rcc.LSEON | lse | rcc.RTCEN

	wait()

	if RCC.BDCR.Bits(rcc.LSEON|rcc.RTCSEL|rcc.RTCEN) != bdcrcfg {
		fmt.Println("Configuring backup domain...")
		RCC.PWREN().Set()
		_ = RCC.PWREN().Load()
		PWR.DBP().Set()
		RCC.BDRST().Set()
		RCC.BDRST().Clear()
		RCC.BDCR.StoreBits(rcc.LSEON|rcc.RTCSEL|rcc.RTCEN, bdcrcfg)
		for RCC.LSERDY().Load() == 0 {
		}
		PWR.DBP().Clear()
		RCC.PWREN().Clear()
		fmt.Println("Done.")
	}
	if RTC.INITS().Load() == 0 {
		fmt.Println("RTC not initialized. Initializing...")
		RCC.PWREN().Set()
		_ = RCC.PWREN().Load()
		PWR.DBP().Set()
		RTC.WPR.Store(0xca)
		RTC.WPR.Store(0x53)
		RTC.INIT().Set()
		for RTC.INITF().Load() == 0 {
		}
		RTC.PRER.Store((4-1)<<16 + (8192 - 1))
		RTC.PRER.Store((4-1)<<16 + (8192 - 1))
		RTC.DR.Store(0x151215 + 2<<13)
		RTC.TR.Store(0x214540)
		RTC.INIT().Clear()
		RTC.WPR.Store(0xff)
		PWR.DBP().Clear()
		RCC.PWREN().Clear()
		fmt.Println("Done.")
	}

	for {
		LEDport.SetPin(Green)
		wait()
		LEDport.ClearPin(Green)
		wait()
		ss := RTC.SSR.Load()
		hhmmss := RTC.TR.Load()
		yymmdd := RTC.DR.Load() &^ (7 << 13)
		ms := 1000 * (8192 - ss) / 8192
		fmt.Printf("20%06x %06x.%03d\n", yymmdd, hhmmss, ms)
	}
}
