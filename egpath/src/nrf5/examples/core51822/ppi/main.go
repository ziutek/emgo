// This example show how to use Programmable Peripheral Interconnect (PPI)
// peripheral. CPU only configures input and output pins, a timer, PPI channels
// and groups (see init function). After that it falls asleep and all events
// (from timer and key) are handled by PPI. Timer is connected to LED0 and
// toggles it using GPIOTE. KEY1 when pressed temporarily disables connection
// between timer and LED0 using PPI channel group.
package main

import (
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

var (
	leds [5]gpio.Pin
	key  gpio.Pin
)

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// Allocate pins (always do it in one place to avoid conflicts).
	p0 := gpio.P0
	key = p0.Pin(16) // KEY1
	for i := range leds {
		leds[i] = p0.Pin(18 + i)
	}

	// Configure pins.
	for _, led := range leds {
		led.Setup(gpio.ModeOut)
	}
	key.Setup(gpio.ModeIn | gpio.PullUp)

	t := timer.TIMER1
	t.StorePRESCALER(6) // 250 kHz
	t.StoreCC(0, 0)
	t.Task(timer.START).Trigger()

	led0 := gpiote.Chan(0)
	led0.Setup(leds[0], gpiote.ModeTask|gpiote.PolarityToggle)
	keyL := gpiote.Chan(2)
	keyL.Setup(key, gpiote.ModeEvent|gpiote.PolarityHiToLo)
	keyH := gpiote.Chan(3)
	keyH.Setup(key, gpiote.ModeEvent|gpiote.PolarityLoToHi)

	pc := ppi.Chan(0)
	pc.SetEEP(t.Event(timer.COMPARE(0)))
	pc.SetTEP(led0.OUT().Task())
	pc.Enable()

	pg := ppi.Group(0)
	pg.SetChannels(pc.Mask())

	pc = ppi.Chan(1)
	pc.SetEEP(keyL.IN().Event())
	pc.SetTEP(pg.DIS().Task())
	pc.Enable()
	pc = ppi.Chan(2)
	pc.SetEEP(keyH.IN().Event())
	pc.SetTEP(pg.EN().Task())
	pc.Enable()
}

func main() {
	for {
		delay.Millisec(1e6)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
