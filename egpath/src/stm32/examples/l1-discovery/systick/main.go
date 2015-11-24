// This example shows how to manually setup interrupt table when you don't use
// runtime initialisation (MaxTasks == 0) and how to write purely interrupt
// driven application.
package main

import (
	"sync/fence"

	"arch/cortexm/sleep"
	"arch/cortexm/systick"

	"stm32/l1/gpio"
	"stm32/l1/irqs"
	"stm32/l1/periph"
	"stm32/l1/setup"
)

var LED = gpio.B

const (
	Blue  = 6
	Green = 7
)

func defaultHandler() {
	for {
	}
}

var (
	cnt   int
	ledup = true
)

func isr() {
	if ledup {
		LED.SetPin(Blue)
	} else {
		LED.ClearPin(Blue)
	}
	ledup = !ledup
}

//c:const
//c:__attribute__((section(".InterruptVectors")))
var IRQs = [...]func(){
	irqs.Tim2: isr,
}

func main() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)

	LED.SetMode(Blue, gpio.Out)
	LED.SetMode(Green, gpio.Out)

	tenms := systick.CALIB.Bits(systick.TENMS)
	tenms *= 10 // stm32l1 returns value for 1 ms not for 10ms.
	systick.RVR.StoreBits(systick.RELOAD, tenms*100)
	systick.CVR.StoreBits(systick.CURRENT, 0)
	systick.CSR.SetBits(systick.ENABLE | systick.TICKINT)

	// Sleep forever.
	sleep.EnableSleepOnExit()
	fence.Sync() // not necessary on Cortex-M0,M3,M4
	sleep.WFI()

	// Execution should never reach there so the green LED
	// should never light up.
	LED.SetPin(Green)
}
