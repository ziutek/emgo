package main

import (
	//"debug/semihosting"
	"delay"
	"rtos"
	"sync/fence"
	"unsafe"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/radio"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var (
	leds       [5]gpio.Pin
	radioEvent rtos.EventFlag
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	for i := range leds {
		led := gpio.P0.Pin(18 + i)
		led.Setup(&gpio.Config{Mode: gpio.Out})
		leds[i] = led
	}

	r := radio.RADIO
	r.SetPCNF0(0)
	r.SetPCNF1(radio.MaxLen(2) | radio.StatLen(2) | radio.BALen(2) | radio.MSBFirst)
	r.SetCRCCNF(2, false)
	r.SetCRCPOLY(1<<16 | 1<<12 | 1<<5 | 1)
	r.SetCRCINIT(0xFFFF)
	r.SetBASE0(0xE7E70000) // Reversed 0xE7E70000.
	r.SetPREFIX0(0xE7)     // Reversed 0xE7.
	r.SetTXADDRESS(0)
	r.SetMODE(radio.NRF_250K)
	r.SetFREQUENCY(radio.Channel(50))
	r.SetSHORTS(radio.READY_START | radio.END_DISABLE)
	rtos.IRQ(r.IRQ()).Enable()

	//f, err := semihosting.OpenFile(":tt", semihosting.W)
	//for err != nil {
	//}
	//fmt.DefaultWriter = f
}

func main() {
	var data [2]int8

	r := radio.RADIO
	r.SetPACKETPTR(unsafe.Pointer(&data[0]))

	leds[0].Set()
	dir := 1
	n := 0
	for {
		data[0] = int8(n)
		data[1] = int8(n)

		disev := r.Event(radio.DISABLED)
		disev.Clear()
		disev.EnableIRQ()
		radioEvent.Reset(0)
		fence.W()
		r.Task(radio.TXEN).Trigger()
		radioEvent.Wait(1, 0)
		switch n {
		case 64:
			dir = -1
			leds[0].Clear()
		case -64:
			dir = 1
			leds[0].Set()
		}
		leds[1].Store(n & 1)
		n += dir
		delay.Millisec(25)
	}
}

func radioISR() {
	radio.RADIO.DisableIRQ(0xFFFFFFFF)
	radioEvent.Signal(1)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:  rtcst.ISR,
	irq.RADIO: radioISR,
}
