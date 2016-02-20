// This example shows how to write bare metal application that does not use
// tasker and rely at limited runtime initialisation (MaxTasks == 0).
// Additionaly this is example of a purly interrupt driven application
package main

import (
	"arch/cortexm"
	"arch/cortexm/scb"
	"arch/cortexm/systick"

	"stm32/hal/gpio"
	"stm32/hal/system"
)

const (
	Red  = gpio.Pin14
	Blue = gpio.Pin15
)

var (
	leds  gpio.Port
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
	system.Setup168(8)

	gpio.D.EnableClock(false)
	leds = gpio.D

	cfg := &gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Blue|Red, cfg)

	st := systick.SYSTICK
	onesec := systick.RVR_Bits(system.AHB.Clock() / 8)
	st.RELOAD().Store(onesec/2 - 1) // Period 0.5 s.
	st.CURRENT().Clear()
	st.CSR.SetBits(systick.ENABLE | systick.TICKINT)

	// Sleep forever.
	scb.SCB.SLEEPONEXIT().Set()
	cortexm.DSB() // not necessary on Cortex-M0,M3,M4
	cortexm.WFI()

	// Execution should never reach there so the red LED
	// should never light up.
	leds.SetPins(Red)
}

//emgo:const
//c:__attribute__((section(".SysTick")))
var SysTickVector = sysTickHandler
