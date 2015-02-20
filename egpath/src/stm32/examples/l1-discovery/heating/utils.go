package main

import (
	"delay"
)

type LED uint

func (led LED) On() {
	ledsPort.SetBit(uint(led))
}

func (led LED) Off() {
	ledsPort.ClearBit(uint(led))
}

func (led LED) Blink(ms int) {
	led.On()
	if ms > 0 {
		delay.Millisec(ms)
	} else {
		delay.Loop(-ms)
	}
	led.Off()
}