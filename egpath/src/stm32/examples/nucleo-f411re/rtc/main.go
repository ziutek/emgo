package main

import (
	"delay"
	"fmt"

	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f411xe/gpio"
	"stm32/f411xe/pwr"
	"stm32/f411xe/rcc"
	"stm32/f411xe/rtc"
)

var LED = gpio.GPIOA

const Green = 1 << 5

func init() {
	setup.Performance84(8)
	periph.AHB1ClockEnable(periph.GPIOA)
	periph.AHB1Reset(periph.GPIOA)
	gpio.GPIOA.MODER5().Store(gpio.MODER_OUT * gpio.MODER5_0) // Green LED
}

func wait() {
	delay.Millisec(250)
}

func main() {
	delay.Millisec(500)

	PWR := pwr.PWR
	RCC := rcc.RCC
	RTC := rtc.RTC

	const bdcrcfg = rcc.LSEON | rcc.RTCSEL_LSE | rcc.RTCEN
	if RCC.BDCR.Bits(rcc.LSEON|rcc.RTCSEL|rcc.RTCEN) != bdcrcfg {
		fmt.Println("Configuring backup domain...")
		periph.APB1ClockEnable(periph.PWR)
		PWR.DBP().Set()
		RCC.BDRST().Set()
		RCC.BDRST().Clear()
		RCC.BDCR.StoreBits(rcc.LSEON|rcc.RTCSEL|rcc.RTCEN, bdcrcfg)
		for RCC.LSERDY().Load() == 0 {
		}
		PWR.DBP().Clear()
		periph.APB1ClockDisable(periph.PWR)
		fmt.Println("Done.")
	}
	if RTC.INITS().Load() == 0 {
		fmt.Println("RTC not initialized. Initializing...")
		periph.APB1ClockEnable(periph.PWR)
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
		periph.APB1ClockDisable(periph.PWR)
		fmt.Println("Done.")
	}

	for {
		LED.BSRRL.Store(Green)
		wait()
		LED.BSRRH.Store(Green)
		wait()
		ss := RTC.SSR.Load()
		hhmmss := RTC.TR.Load()
		yymmdd := RTC.DR.Load() &^ (7 << 13)
		ms := 1000 * (8192 - ss) / 8192
		fmt.Printf("20%06x %06x.%03d\n", yymmdd, hhmmss, ms)
	}
}
