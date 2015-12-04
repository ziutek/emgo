package main

import (
	"delay"
	"rtos"
	"sync/fence"

	"arch/cortexm"
	"arch/cortexm/scb"

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

const period = 2 * 32768 // 2s

func init() {
	setup.Clocks(clock.Xtal, clock.Xtal, true)

	for _, led := range leds {
		p0.SetMode(int(led), gpio.Out)
	}

	t0.SetPrescaler(8) // 62500 Hz
	t0.SetCC(1, 65526/2)
	t0.Event(timer.COMPARE0).EnableInt()
	t0.Event(timer.COMPARE1).EnableInt()
	rtos.IRQ(irqs.Timer0).Enable()
	t0.Task(timer.START).Trig()

	rtc0.SetPrescaler(0) // 32768 Hz
	rtc0.Event(rtc.COMPARE1).EnableInt()
	rtos.IRQ(irqs.RTC0).Enable()
	rtc0.SetCC(1, period)
	rtc0.Task(rtc.START).Trig()
}

func blink(led byte, dly int) {
	p0.SetPin(int(led))
	delay.Loop(dly)
	p0.ClearPin(int(led))
}

func timerISR() {
	if e := t0.Event(timer.COMPARE0); e.Happened() {
		e.Clear()
		blink(leds[3], 1e3)
	}
	if e := t0.Event(timer.COMPARE1); e.Happened() {
		e.Clear()
		blink(leds[4], 1e3)
	}
}

func rtcISR() {
	rtc0.Event(rtc.COMPARE1).Clear()
	rtc0.SetCC(1, rtc0.CC(1)+period)
	blink(leds[2], 1e3)
}

var ISRs = [...]func(){
	irqs.Timer0: timerISR,
	irqs.RTC0:   rtcISR,
} //c:__attribute__((section(".InterruptVectors")))

func main() {
	// Sleep forever.
	scb.SLEEPONEXIT.Set()
	fence.Sync() // not necessary on Cortex-M0,M3,M4
	cortexm.WFI()
	// Execution should never reach there so LED0 should never light up.
	p0.SetPin(int(leds[0]))
}
