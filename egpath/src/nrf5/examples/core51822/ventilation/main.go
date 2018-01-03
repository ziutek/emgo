package main

import (
	"delay"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

// Two BW428G-E4 four digits, LED displays (common cathode).
type Display struct {
	dig [8]gpio.Pin // 0-3 top display, 4-7 bottom display.
	seg [8]gpio.Pin // A B C D E F G :
}

var disp Display

func init() {
	// Initialize system and runtime.
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	// Allocate pins (always do it in one place to avoid conflicts).

	p0 := gpio.P0
	disp.seg[5] = p0.Pin(0)  // F
	disp.dig[5] = p0.Pin(1)  // Bottom 1
	disp.dig[6] = p0.Pin(2)  // Bottom 2
	disp.seg[3] = p0.Pin(3)  // D
	disp.seg[4] = p0.Pin(9)  // E
	disp.dig[7] = p0.Pin(11) // Bottom 3
	disp.dig[2] = p0.Pin(15) // Top 2
	disp.seg[6] = p0.Pin(17) // 6
	disp.seg[1] = p0.Pin(22) // B
	disp.seg[2] = p0.Pin(23) // C
	disp.seg[7] = p0.Pin(21) // :
	disp.dig[0] = p0.Pin(24) // Top 0
	disp.dig[3] = p0.Pin(25) // Top 3
	disp.seg[0] = p0.Pin(28) // A
	disp.dig[1] = p0.Pin(29) // Top 1
	disp.dig[4] = p0.Pin(30) // Bottom 0

	// Configure pins.

	for _, pin := range disp.dig {
		if !pin.IsValid() {
			continue
		}
		// Drive digits with higd drive, open collector.
		pin.Set()
		pin.Setup(gpio.ModeOut | gpio.DriveH0D1)
	}
	for _, pin := range disp.seg {
		if !pin.IsValid() {
			continue
		}
		pin.Setup(gpio.ModeOut | gpio.DriveD0H1)
	}
}

func wait() {
	delay.Millisec(1000)
}

func main() {
	disp.seg[0].Set()
	disp.seg[1].Set()
	disp.seg[2].Set()
	disp.seg[3].Set()
	disp.seg[4].Set()
	disp.seg[5].Set()
	disp.seg[6].Set()
	disp.seg[7].Set()
	for {
		disp.dig[7].Set()
		disp.dig[0].Clear()
		wait()
		disp.dig[0].Set()
		disp.dig[1].Clear()
		wait()
		disp.dig[1].Set()
		disp.dig[2].Clear()
		wait()
		disp.dig[2].Set()
		disp.dig[3].Clear()
		wait()
		disp.dig[3].Set()
		disp.dig[4].Clear()
		wait()
		disp.dig[4].Set()
		disp.dig[5].Clear()
		wait()
		disp.dig[5].Set()
		disp.dig[6].Clear()
		wait()
		disp.dig[6].Set()
		disp.dig[7].Clear()
		wait()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
