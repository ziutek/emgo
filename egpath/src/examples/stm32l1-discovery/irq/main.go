package main

import (
	"cortexm/irq"
	"cortexm/sleep"
	"cortexm/startup"
	"cortexm/systick"

	"stm32/l1/gpio"
	"stm32/l1/periph"
	"stm32/l1/setup"
)

// STM32L1-Discovery LEDs
var LEDs = gpio.B

const (
	Blue  = 6
	Green = 7
)

var vt = irq.SysTable{
	Reset:     startup.Start,
	NMI:       defaultHandler,
	HardFault: defaultHandler,
	SysTick:   sysTickHandler,
}

func defaultHandler() {
	for {
	}
}

var (
	cnt int
	led = true
)

func sysTickHandler() {
	if cnt++; cnt < 100 {
		return
	}
	cnt = 0

	if led {
		LEDs.SetBit(Blue)
	} else {
		LEDs.ClearBit(Blue)
	}
	led = !led
}

func main() {
	setup.Performance(0)

	periph.AHBClockEnable(periph.GPIOB)
	periph.AHBReset(periph.GPIOB)

	LEDs.SetMode(Blue, gpio.Out)
	LEDs.SetMode(Green, gpio.Out)

	irq.SetActiveTable(vt.Slice())

	_, _, tenms := systick.Calib()
	tenms *= 10 // stm32l1 returns value for 1 ms.
	systick.SetReload(tenms)
	systick.SetFlags(systick.Enable | systick.TickInt)

	sleep.EnableSleepOnExit()
	sleep.WFI()
	
	// Execution should never reach there so Green LED
	// should never light up.
	LEDs.SetBit(Green)
}
