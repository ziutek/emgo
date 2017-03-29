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

	gtec := gpiote.Chan(0)
	cfg := gpiote.Config{
		Mode:     gpiote.Task,
		Polarity: gpiote.Toggle,
	}
	gtec.Setup(gpio.P0.Pin(18), cfg)

	t := timer.TIMER0
	t.StorePRESCALER(8) // 16 MHz / 2^8 = 62500 Hz
	t.StoreCC(1, 62500/4)
	t.StoreSHORTS(timer.COMPARE1_CLEAR)

	ppic := ppi.Chan(0)
	ppic.SetEEP(t.Event(timer.COMPARE1))
	ppic.SetTEP(gtec.OUT())
	ppic.Enable()

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
