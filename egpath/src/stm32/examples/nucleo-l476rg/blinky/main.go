package main

import (
	"delay"
	"fmt"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var led gpio.Pin

func init() {
	//system.SetupPLL(-48, 6, 20, 0, 0, 2) // 80 MHz (max. for voltage Range 1).
	//system.SetupPLL(-4, 1, 26, 0, 0, 4) // 26 MHz (max. for voltage Range 2).
	system.SetupMSI(100) // Lowest possible frequency.
	systick.Setup(1e7) // Typical 2e6 ns is to low for 100 kHz SysClk.

	gpio.A.EnableClock(false)
	led = gpio.A.Pin(5)

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	led.Setup(&cfg)
}

func main() {
	delay.Millisec(500)
	buses := []system.Bus{
		system.Core,
		system.AHB,
		system.APB1,
		system.APB2,
	}
	fmt.Printf("\r\n")
	for _, bus := range buses {
		fmt.Printf("%4s: %9d Hz\r\n", bus, bus.Clock())
	}
	for {
		led.Set()
		delay.Millisec(100)
		led.Clear()
		delay.Millisec(1900)
	}
}
