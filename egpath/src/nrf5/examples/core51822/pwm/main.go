package main

import (
	//"debug/semihosting"
	"delay"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/gpiote"
	"nrf5/hal/irq"
	"nrf5/hal/ppi"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/timer"
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	te := gpiote.GPIOTE
	cfg := gpiote.Config{
		Pin:      gpio.P0.Pin(18), // LED0
		Mode:     gpiote.Task,
		Polarity: gpiote.Toggle,
	}
	te.StoreCONFIG(0, cfg)

	t := timer.TIMER0
	t.StorePRESCALER(8) // 16 MHz / 2^8 = 62500 Hz
	t.StoreCC(1, 62500/4)
	t.StoreSHORTS(timer.COMPARE1_CLEAR)

	c := ppi.Chan(0)
	c.SetEEP(t.Event(timer.COMPARE1))
	c.SetTEP(te.OUT(0))
	c.Enable()

	t.Task(timer.START).Trigger()
}

func main() {
	for {
		delay.Millisec(1000)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
