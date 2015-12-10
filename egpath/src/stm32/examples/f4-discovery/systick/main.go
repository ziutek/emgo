// This example shows how to manually setup interrupt table when you don't use
// runtime initialisation (MaxTasks == 0) and how to write purely interrupt
// driven application.
package main

import (
	"arch/cortexm"
	"arch/cortexm/scb"
	"arch/cortexm/systick"

	"stm32/f4/gpio"
	"stm32/f4/periph"
	"stm32/f4/setup"
)

var LED = gpio.D

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

var (
	cnt   int
	ledup = true
)

func sysTickHandler() {
	if ledup {
		LED.SetPin(Blue)
	} else {
		LED.ClearPin(Blue)
	}
	ledup = !ledup
}

//c:__attribute__((section(".SysTick")))
var SysTickVector = sysTickHandler

func main() {
	setup.Performance168(8)

	periph.AHB1ClockEnable(periph.GPIOD)
	periph.AHB1Reset(periph.GPIOD)

	LED.SetMode(Blue, gpio.Out)
	LED.SetMode(Red, gpio.Out)

	tenms := systick.TENMS.Load()
	tenms *= 10
	systick.RELOAD.Store(tenms * 100)
	systick.CURRENT.Clear()
	(systick.ENABLE | systick.TICKINT).Set()

	// Sleep forever.
	scb.SLEEPONEXIT.Set()
	cortexm.DSB() // not necessary on Cortex-M0,M3,M4
	cortexm.WFI()

	// Execution should never reach there so the red LED
	// should never light up.
	LED.SetPin(Red)
}
