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
	clkm := clock.Mgmt
	clkm.TrigTask(clock.HFCLKSTART)
	for {
		src, run := clkm.HFCLKStat()
		if run && src == clock.Xtal {
			break
		}
	}
	clkm.SetLFCLKSrc(clock.Xtal)
	clkm.TrigTask(clock.LFCLKSTART)
	for {
		src, run := clkm.LFCLKStat()
		if run && src == clock.Xtal {
			break
		}
	}

	for _, led := range leds {
		p0.SetMode(int(led), gpio.Out)
	}

	t0.SetPrescaler(8) // 62500 Hz
	t0.SetCC(1, 65526/2)
	t0.EnableInt(timer.COMPARE0)
	t0.EnableInt(timer.COMPARE1)
	rtos.IRQ(t0.IRQ()).Enable()
	t0.TrigTask(timer.START)

	rtc0.SetPrescaler(1<<12 - 1) // 8 Hz
	rtc0.EnableEvent(rtc.TICK)
	rtc0.EnableInt(rtc.TICK)
	rtos.IRQ(rtc0.IRQ()).Enable()
	rtc0.TrigTask(rtc.START)

}

func blink(led byte, dly int) {
	p0.SetPin(int(led))
	delay.Loop(dly)
	p0.ClearPin(int(led))
}

func timerISR() {
	switch {
	case t0.Event(timer.COMPARE0):
		t0.ClearEvent(timer.COMPARE0)
		blink(leds[3], 1e3)
	case t0.Event(timer.COMPARE1):
		t0.ClearEvent(timer.COMPARE1)
		blink(leds[4], 1e3)
	}
}

func rtcISR() {
	rtc0.ClearEvent(rtc.TICK)
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
