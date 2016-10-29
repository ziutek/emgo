package sysclk

import (
	"nrf51/irq"
	"nrf51/rtc"
)

const freqHz = 32768

var rtc0 = rtc.RTC0

func init() {
	rtc0.SetPrescaler(0) // 32768 Hz
	rtc0.Event(rtc.COMPARE1).EnableInt()
	irq.RTC0.Enable()
	rtc0.Task(rtc.START).Trig()

}

func isr() {
	if e := rtc0.Event(rtc.OVRFLW); e.Happened() {
		e.Clear()

	}
	if e := rtc0.Event(rtc.COMPARE0); e.Happened() {
		e.Clear()

	}
}

//emgo:const
//c:__attribute__((section(".InterruptVectors")))
var ISRs = [...]func(){
	irq.RTC0: isr,
}
