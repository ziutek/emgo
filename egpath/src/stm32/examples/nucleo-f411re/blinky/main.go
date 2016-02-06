package main

import (
	"delay"

	"arch/cortexm/bitband"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var led bitband.Bit

func init() {
	system.Setup96(8)
	systick.Setup()

	port, pin := gpio.A, 5
	led = port.OutPin(pin)

	port.EnableClock(false)

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	gpio.A.SetupPin(pin, &cfg)
}

func wait() {
	delay.Millisec(500)
}

func main() {
	for {
		led.Set()
		wait()
		led.Clear()
		wait()
	}
}
