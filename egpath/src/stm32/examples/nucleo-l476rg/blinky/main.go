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
	system.Setup(-48, 6, 20, 0, 0, 2) // 80 MHz (fastest for voltage Range 1).
	//system.Setup(-4, 1, 26, 0, 0, 4) // 26 MHz (fastest for voltage Range 2).
	systick.Setup()

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
		delay.Millisec(50)
		led.Clear()
		delay.Millisec(950)
	}
}
