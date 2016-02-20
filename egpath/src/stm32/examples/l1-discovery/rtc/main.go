package main

import (
	"delay"
	"fmt"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/pwr"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/rtc"
)

var leds gpio.Port

const (
	Green = gpio.Pin7
)

func init() {
	system.Setup32(0)
	systick.Setup()

	gpio.B.EnableClock(false)
	leds = gpio.B

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green, &cfg)
}

func main() {
	delay.Millisec(500)

	RTC := rtc.RTC
	RCC := rcc.RCC
	PWR := pwr.PWR

	const (
		mask = rcc.LSEON | rcc.RTCSEL | rcc.RTCEN
		cfg  = rcc.LSEON | rcc.RTCSEL_LSE | rcc.RTCEN
	)
	if RCC.CSR.Bits(mask) != cfg {
		fmt.Println("Configuring backup domain...")
		RCC.PWREN().Set()
		_ = RCC.PWREN().Load()
		PWR.DBP().Set()
		RCC.RTCRST().Set()
		RCC.RTCRST().Clear()
		RCC.CSR.StoreBits(mask, cfg)
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
		const prer = (2-1)<<16 + (1 - 1)
		RTC.PRER.Store(prer)
		RTC.PRER.Store(prer)
		fmt.Printf("%x\n", RTC.PRER.Load())
		RTC.DR.Store(0x151215 + 2<<13)
		RTC.TR.Store(0x214540)
		RTC.INIT().Clear()
		RTC.WPR.Store(0xff)
		PWR.DBP().Clear()
		RCC.PWREN().Clear()
		fmt.Println("Done.")
	}

	for {
		leds.SetPins(Green)
		delay.Millisec(250)
		leds.ClearPins(Green)
		delay.Millisec(250)
		hhmmss := RTC.TR.Load()
		dr := RTC.DR.Load()
		yymmdd := dr &^ (7 << 13)
		w := (dr >> 13) & 7
		fmt.Printf("%06x %d %06x\n", yymmdd, w, hhmmss)
	}
}
