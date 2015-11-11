// This example shows how to manually setup interrupt table when you don't use
// runtime initialisation (MaxTasks == 0) and how to write purely interrupt
// driven application.
package main

import (
	"sync/barrier"

	"arch/cortexm/exce"
	"arch/cortexm/sleep"
	"arch/cortexm/systick"

	"stm32/l1/gpio"
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

func sysTickHandler() {
	if ledup {
		LED.SetPin(Blue)
	} else {
		LED.ClearPin(Blue)
	}
	ledup = !ledup
}

func main() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)

	LED.SetMode(Blue, gpio.Out)
	LED.SetMode(Green, gpio.Out)

	vt := exce.NewTable(16)
	vt[exce.NMI] = exce.VectorFor(defaultHandler)
	vt[exce.HardFault] = exce.VectorFor(defaultHandler)
	vt[exce.SysTick] = exce.VectorFor(sysTickHandler)
	exce.UseTable(vt)

	_, _, tenms := systick.Calib()
	tenms *= 10 // stm32l1 returns value for 1 ms not for 10ms.
	systick.SetReload(tenms * 100)
	systick.SetFlags(systick.Enable | systick.TickInt)

	// Sleep forever.
	sleep.EnableSleepOnExit()
	barrier.Sync() // not necessary on Cortex-M0,M3,M4
	sleep.WFI()

	// Execution should never reach there so the green LED
	// should never light up.
	LED.SetPin(Green)
}
