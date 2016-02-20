package main

import (
	"delay"
)

type LED int

func (led LED) On() {
	ledsPort.SetPin(int(led))
}

func (led LED) Off() {
	ledsPort.ClearPin(int(led))
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