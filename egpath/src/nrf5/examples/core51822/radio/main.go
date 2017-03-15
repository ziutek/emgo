package main

import (
	"debug/semihosting"
	"fmt"
	"unsafe"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/radio"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var leds [5]gpio.Pin

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	for i := range leds {
		led := gpio.P0.Pin(18 + i)
		led.Setup(&gpio.Config{Mode: gpio.Out})
		leds[i] = led
	}

	f, err := semihosting.OpenFile(":tt", semihosting.W)
	for err != nil {
	}
	fmt.DefaultWriter = f
}

func main() {

	r := radio.RADIO

	r.SetPCNF0(radio.MakePktConf0(8, 0, 0, 0, false))
	r.SetPCNF1(radio.MakePktConf1(2, 2, 2, true, false))
	r.SetCRCCNF(2, false)
	r.SetCRCPOLY(1<<16 | 1<<12 | 1<<5 | 1)
	r.SetCRCINIT(0xFFFF)
	r.SetBASE0(0xAC0F) // Reversed 0xF035.
	r.SetPREFIX0(0xEE) // Reversed 0x77.
	r.SetTXADDRESS(0)
	r.SetMODE(radio.NRF_250K)
	r.SetFREQUENCY(radio.MakeFreq(2450, false))

	payload := uint16(0x1234)
	r.SetPACKETPTR(uintptr(unsafe.Pointer(&payload)))

	fmt.Println(r.STATE())

	leds[0].Set()

	r.Event(radio.READY).Clear()
	r.Task(radio.TXEN).Trigger()
	for !r.Event(radio.READY).IsSet() {
		fmt.Println(r.STATE())
	}
	fmt.Println(r.STATE())

	leds[1].Set()

	r.Event(radio.END).Clear()
	r.Task(radio.START).Trigger()
	for !r.Event(radio.END).IsSet() {
		fmt.Println(r.STATE())
	}
	fmt.Println(r.STATE())

	leds[2].Set()

	r.Event(radio.DISABLED).Clear()
	r.Task(radio.DISABLE).Trigger()
	for !r.Event(radio.DISABLED).IsSet() {
		fmt.Println(r.STATE())
	}
	fmt.Println(r.STATE())

	leds[3].Set()

}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
