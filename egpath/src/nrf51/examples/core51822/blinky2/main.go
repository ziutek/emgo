package main

import (
	"delay"
	"rtos"
	"syscall"

	"nrf51/clock"
	"nrf51/gpio"
	"nrf51/rtc"
	"nrf51/setup"
)

//emgo:const
var leds = [...]byte{18, 19, 20, 21, 22}

var (
	p0   = gpio.P0
	rtc0 = rtc.RTC0
)

func init() {
	setup.Clocks(clock.Xtal, clock.Xtal, true)

	for _, led := range leds {
		p0.SetMode(int(led), gpio.Out)
	}

	rtc0.SetPRESCALER(0) // 32768 Hz
	rtc0.TASK(rtc.START).Trigger()
}

func blink(led byte, dly int) {
	p0.SetPin(int(led))
	delay.Loop(dly)
	p0.ClearPin(int(led))
	delay.Loop(dly)
}

func uptime() int64 {
	return int64(rtc0.COUNTER()) * 1e9 / 32768
}

func main() {
	syscall.SetSysTimer(uptime, nil)
	for t := int64(0); ; t += 1e9 {
		for rtos.Nanosec() < t {
		}
		blink(leds[2], 1)
	}
}
