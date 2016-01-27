// This example shows how to manually setup interrupt table when you don't use
// runtime initialisation (MaxTasks == 0) and how to write purely interrupt
// driven application.
package main

import (
	"arch/cortexm"
	"arch/cortexm/scb"
	"arch/cortexm/systick"

	"stm32/hal/gpio"
	"stm32/hal/system"
)

var LED *gpio.Port

const (
	Red  = gpio.Pin14
	Blue = gpio.Pin15
)

var ledup = true

func sysTickHandler() {
	if ledup {
		LED.SetPins(Blue)
	} else {
		LED.ClearPins(Blue)
	}
	ledup = !ledup
}

//c:__attribute__((section(".SysTick")))
var SysTickVector = sysTickHandler

func main() {
	system.Setup168(8)

	gpio.D.EnableClock(false)
	LED = gpio.D

	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	LED.Setup(Blue|Red, cfg)

	st := systick.SYSTICK
	onesec := systick.RVR_Bits(system.AHBClk / 8)
	st.RELOAD().Store(onesec/2 - 1) // Period 0.5 s.
	st.CURRENT().Clear()
	st.CSR.SetBits(systick.ENABLE | systick.TICKINT)

	// Sleep forever.
	scb.SLEEPONEXIT.Set()
	cortexm.DSB() // not necessary on Cortex-M0,M3,M4
	cortexm.WFI()

	// Execution should never reach there so the red LED
	// should never light up.
	LED.SetPins(Red)
}
