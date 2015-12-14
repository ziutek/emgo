package main

import (
	"delay"
	"fmt"

	"stm32/f4/gpio"
	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f411/pwr"
	"stm32/f411/rcc"
	"stm32/rtc"
)

var LED = gpio.A

const (
	Green = 5
)

func init() {
	setup.Performance84(8)
	periph.AHB1ClockEnable(periph.GPIOA)
	periph.AHB1Reset(periph.GPIOA)
	LED.SetMode(Green, gpio.Out)
}

func wait() {
	delay.Millisec(50)
}

func main() {
	delay.Millisec(1000)

	const bdcrcfg = rcc.LSEON | rcc.RTCSEL_LSE | rcc.RTCEN
	if rcc.BDCR_Load()&(rcc.LSEON|rcc.RTCSEL|rcc.RTCEN) != bdcrcfg {
		fmt.Println("Configuring backup domain...")
		periph.APB1ClockEnable(periph.PWR)
		pwr.DBP.Set()
		rcc.BDRST.Set()
		rcc.BDCR_Store(bdcrcfg)
		for rcc.LSERDY.Load() == 0 {
		}
		pwr.DBP.Clear()
		periph.APB1ClockDisable(periph.PWR)
		fmt.Println("Done.")
	}
	if rtc.INITS.Load() == 0 {
		fmt.Println("RTC not initialized. Initializing...")
		periph.APB1ClockEnable(periph.PWR)
		pwr.DBP.Set()
		rtc.WPR_Store(0xca)
		rtc.WPR_Store(0x53)
		rtc.INIT.Set()
		for rtc.INITF.Load() == 0 {
		}
		rtc.PREDIV_S.Store(32767)
		rtc.PREDIV_A.Store(0)
		rtc.DR_Store(0x151214 + 1<<13)
		rtc.TR_Store(0x013514)
		rtc.INIT.Clear()
		rtc.WPR_Store(0xff)
		pwr.DBP.Clear()
		periph.APB1ClockDisable(periph.PWR)
		fmt.Println("Done.")
	}

	for {
		LED.SetPin(Green)
		wait()
		LED.ClearPin(Green)
		wait()
		ss := uint32(rtc.SSR_Load())
		hhmmss := rtc.TR_Load()
		yymmdd := rtc.DR_Load() &^ (1 << 13)
		ms := 1000 * (32767 - ss) / 32768
		fmt.Printf("%06x %06x.%03d\n", yymmdd, hhmmss, ms)
		
		// Sprawdzić jaki kod generowany jest dla maski 0xffffffff
	}
}


/*
Z uintami:
   text	   data	    bss	    dec	    hex	filename
  32060	    312	  19572	  51944	   cae8	cortexm4f.elf
Bez uintow:
   text	   data	    bss	    dec	    hex	filename
  32860	    416	  19792	  53068	   cf4c	cortexm4f.elf

*/