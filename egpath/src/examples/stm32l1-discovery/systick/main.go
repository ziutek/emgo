package main

import (
	"sync"

	"cortexm/irq"
	"cortexm/sleep"
	"cortexm/startup"
	"cortexm/systick"

	"stm32/l1/gpio"
	"stm32/l1/periph"
	"stm32/l1/setup"
)

// STM32L1-Discovery LEDs
var LED = gpio.B

const (
	Blue  = 6
	Green = 7
)

var vt = irq.SysTable{
	Reset:     irq.Vector(startup.Start),
	NMI:       irq.Vector(defaultHandler),
	HardFault: irq.Vector(defaultHandler),
	SysTick:   irq.Vector(sysTickHandler),
}

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
		LED.SetBit(Blue)
	} else {
		LED.ClearBit(Blue)
	}
	ledup = !ledup
}

func main() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)

	LED.SetMode(Blue, gpio.Out)
	LED.SetMode(Green, gpio.Out)

	irq.SetActiveTable(vt.Slice())

	_, _, tenms := systick.Calib()
	tenms *= 10 // stm32l1 returns value for 1 ms not for 10ms.
	systick.SetReload(tenms * 100)
	systick.SetFlags(systick.Enable | systick.TickInt)

	// Sleep forever.
	sleep.EnableSleepOnExit()
	sync.Sync() // not necessary on Cortex-M0,M0+,M3,M4
	sleep.WFI()

	// Execution should never reach there so the green LED
	// should never light up.
	LED.SetBit(Green)
}
