package main

import (
	"delay"
	"rtos"

	"arch/cortexm"
	"arch/cortexm/scb"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/timer"
)

const (
	led0 = gpio.Pin18
	led1 = gpio.Pin19
	led2 = gpio.Pin20
	led3 = gpio.Pin21
	led4 = gpio.Pin22
)

//emgo.const
var (
	p0   = gpio.P0
	t0   = timer.TIMER0
	rtc0 = rtc.RTC0
)

const period = 2 * 32768 // 2s

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)

	p0.Setup(led0|led1|led2|led3|led4, &gpio.Config{Mode: gpio.Out})

	t0.StorePRESCALER(8) // 62500 Hz
	t0.StoreCC(1, 65526/2)
	t0.Event(timer.COMPARE0).EnableIRQ()
	t0.Event(timer.COMPARE1).EnableIRQ()
	rtos.IRQ(irq.TIMER0).Enable()
	t0.Task(timer.START).Trigger()

	rtc0.StorePRESCALER(0) // 32768 Hz
	rtc0.Event(rtc.COMPARE1).EnableIRQ()
	rtos.IRQ(irq.RTC0).Enable()
	rtc0.StoreCC(1, period)
	rtc0.Task(rtc.START).Trigger()
}

func blink(led gpio.Pins, dly int) {
	p0.SetPins(led)
	delay.Loop(dly)
	p0.ClearPins(led)
}

func timerISR() {
	if e := t0.Event(timer.COMPARE0); e.IsSet() {
		e.Clear()
		blink(led3, 1e3)
	}
	if e := t0.Event(timer.COMPARE1); e.IsSet() {
		e.Clear()
		blink(led4, 1e3)
	}
}

func rtcISR() {
	rtc0.Event(rtc.COMPARE1).Clear()
	rtc0.StoreCC(1, rtc0.LoadCC(1)+period)
	blink(led2, 1e3)
}

func main() {
	// Sleep forever.
	scb.SCB.SLEEPONEXIT().Set()
	cortexm.DSB() // not necessary on Cortex-M0,M3,M4
	cortexm.WFI()
	// Execution should never reach there so LED0 should never light up.
	p0.SetPins(led0)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.TIMER0: timerISR,
	irq.RTC0:   rtcISR,
}
