package main

import (
	"arch/cortexm/exce"
	"arch/cortexm/sleep"
	"delay"
	"rtos"

	"nrf51/clock"
	"nrf51/gpio"
	"nrf51/irqs"
	"nrf51/rtc"
	"nrf51/setup"
	"nrf51/timer"
)

//c:const
var leds = [...]byte{18, 19, 20, 21, 22}

var (
	p0   = gpio.P0
	t0   = timer.Timer0
	rtc0 = rtc.RTC0
)

func init() {
	setup.Clocks(clock.Xtal, clock.Xtal, true)

	for _, led := range leds {
		p0.SetMode(int(led), gpio.Out)
	}

	t0.SetPrescaler(8) // 62500 Hz
	t0.SetCC(1, 65526/2)
	t0.COMPARE(0).EnableInt()
	t0.COMPARE(1).EnableInt()
	rtos.IRQ(t0.IRQ()).Enable()
	t0.START().Trig()

	rtc0.SetPrescaler(1<<12 - 1) // 8 Hz
	rtc0.TICK().EnableInt()
	rtos.IRQ(rtc0.IRQ()).Enable()
	rtc0.START().Trig()

}

func blink(led byte, dly int) {
	p0.SetPin(int(led))
	delay.Loop(dly)
	p0.ClearPin(int(led))
}

func timerISR() {
	if e := t0.COMPARE(0); e.Happened() {
		e.Clear()
		blink(leds[3], 1e3)
	}
	if e := t0.COMPARE(1); e.Happened() {
		e.Clear()
		blink(leds[4], 1e3)
	}
}

func rtcISR() {
	rtc0.TICK().Clear()
	blink(leds[2], 1e3)
}

//c:const
//c:__attribute__((section(".InterruptVectors")))
var IRQs = [...]func(){
	irqs.Timer0 - exce.IRQ0: timerISR,
	irqs.RTC0 - exce.IRQ0:   rtcISR,
}

func main() {
	for {
		sleep.WFE()
	}
}
