// This example shows how to manually setup interrupt table when you don't use
// runtime initialisation (MaxTasks == 0) and how to write purely interrupt
// driven application.
package main

import (
	"arch/cortexm"
	"arch/cortexm/scb"
	"arch/cortexm/systick"

	"stm32/hal/gpio"
	"stm32/hal/setup"
)

var LED *gpio.Port

const (
	Green  = 12
	Orange = 13
	Red    = 14
	Blue   = 15
)

var ledup = true

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

	gpio.D.EnableClock(false)

	LED = gpio.D
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
