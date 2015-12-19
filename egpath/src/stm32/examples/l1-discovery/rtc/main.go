package main

import (
	"delay"
	"fmt"

	"stm32/l1/gpio"
	"stm32/l1/periph"
	"stm32/l1/setup"
	"stm32/l15x/pwr"
	"stm32/l15x/rcc"
	"stm32/rtc"
)

var LED = gpio.B

const (
	Green = 7
)

func init() {
	setup.Performance(0)
	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)
	LED.SetMode(Green, gpio.Out)
}

func wait() {
	delay.Millisec(250)
}

func main() {
	delay.Millisec(500)

	RTC := rtc.RTC
	RCC := rcc.RCC
	PWR := pwr.PWR

	const csrcfg = rcc.LSEON | rcc.RTCSEL_LSE | rcc.RTCEN
	if RCC.CSR.Bits(rcc.LSEON|rcc.RTCSEL|rcc.RTCEN) != csrcfg {
		fmt.Println("Configuring backup domain...")
		periph.APB1ClockEnable(periph.PWR)
		PWR.DBP().Set()
		RCC.RTCRST().Set()
		RCC.RTCRST().Clear()
		RCC.CSR.StoreBits(rcc.LSEON|rcc.RTCSEL|rcc.RTCEN, csrcfg)
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
		const prer = (2-1)<<16 + (1-1)
		RTC.PRER.Store(prer)
		RTC.PRER.Store(prer)
		fmt.Printf("%x\n", RTC.PRER.Load())
		//RTC.DR.Store(0x151215 + 2<<13)
		//RTC.TR.Store(0x214540)
		RTC.INIT().Clear()
		RTC.WPR.Store(0xff)
		PWR.DBP().Clear()
		periph.APB1ClockDisable(periph.PWR)
		fmt.Println("Done.")
	}

	for {
		LED.SetPin(Green)
		wait()
		LED.ClearPin(Green)
		wait()
		hhmmss := RTC.TR.Load()
		dr := RTC.DR.Load()
		yymmdd := dr &^ (7 << 13)
		w := (dr >> 13) & 7
		fmt.Printf("20%06x %d %06x\n", yymmdd, w, hhmmss)
	}
}
