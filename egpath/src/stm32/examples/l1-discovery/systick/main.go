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

const (
	Blue  = gpio.Pin6
	Green = gpio.Pin7
)

var (
	leds  *gpio.Port
	ledup = true
)

func sysTickHandler() {
	if ledup {
		leds.SetPins(Blue)
	} else {
		leds.ClearPins(Blue)
	}
	ledup = !ledup
}

func main() {
	system.Setup32(0)

	gpio.B.EnableClock(false)
	leds = gpio.B

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Blue|Green, &cfg)

	st := systick.SYSTICK
	onesec := systick.RVR(system.AHB.Clock() / 8)
	st.RELOAD().Store(onesec/2 - 1) // Period 0.5 s.
	st.CURRENT().Clear()
	st.CSR.SetBits(systick.ENABLE | systick.TICKINT)

	// Sleep forever.
	scb.SCB.SLEEPONEXIT().Set()
	cortexm.DSB() // not necessary on Cortex-M0,M3,M4
	cortexm.WFI()

	// Execution should never reach there so the green LED
	// should never light up.
	leds.SetPins(Green)
}

//emgo:const
//c:__attribute__((section(".SysTick")))
var SysTickVector = sysTickHandler
