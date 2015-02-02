package main

import (
	"delay"
)

type LED int

func (led LED) On() {
	ledsPort.SetBit(int(led))
}

func (led LED) Off() {
	ledsPort.ClearBit(int(led))
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

func beep(ms int) {
	for ; ms > 0; ms--{
		buzzPort.SetBit(buzz)
		delay.Millisec(1)
		buzzPort.ClearBit(buzz)
		delay.Millisec(1)
	}
}
