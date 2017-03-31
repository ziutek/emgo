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

const (
	pre  = 1
	freq = 244 // Hz
	max  = 16000000 / (1 << pre) / freq
)

var (
	t    *timer.Periph
	led  gpio.Pin
	gtec gpiote.Chan
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	led = gpio.P0.Pin(18)
	led.Setup(gpio.ModeOut)

	gtec = gpiote.Chan(0)

	t = timer.TIMER1
	t.StorePRESCALER(pre)
	t.StoreCC(0, 1)
	t.StoreCC(1, max)
	t.StoreSHORTS(timer.COMPARE1_CLEAR)

	ppic := ppi.Chan(0)
	ppic.SetEEP(t.Event(timer.COMPARE0))
	ppic.SetTEP(gtec.OUT())
	ppic.Enable()
	ppic = ppi.Chan(1)
	ppic.SetEEP(t.Event(timer.COMPARE1))
	ppic.SetTEP(gtec.OUT())
	ppic.Enable()

	t.Task(timer.START).Trigger()
}

func main() {
	pwmcfg := gpiote.ModeTask |
		gpiote.PolarityToggle |
		gpiote.OutInitHigh
	for {
		for v := uint32(1); v <= max; v *= 2 {
			gtec.Setup(led, 0)
			t.Task(timer.STOP).Trigger()
			t.Task(timer.CLEAR).Trigger()
			t.StoreCC(0, v)
			gtec.Setup(led, pwmcfg)
			t.Task(timer.START).Trigger()

			delay.Millisec(500)
		}
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
