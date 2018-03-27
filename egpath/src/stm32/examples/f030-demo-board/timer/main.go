package main

import (
	"delay"
	"rtos"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/tim"
)

var (
	led1 gpio.Pin
	led2 gpio.Pin
	t    *tim.Periph
	e *tim.CR1
)

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	gpio.A.EnableClock(false)
	led1 = gpio.A.Pin(4)
	led2 = gpio.A.Pin(5)

	cfg := &gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain}
	led1.Setup(cfg)
	led2.Setup(cfg)

	t = tim.TIM3
	pclk := t.Bus().Clock()
	if pclk < system.AHB.Clock() {
		pclk *= 2
	}
	freq := uint(1e3) // Hz
	t.PSC.Store(tim.PSC(pclk / freq))
	t.ARR.Store(250) // ms
	t.DIER.Store(tim.UIE)
	t.CR1.Store(tim.CEN)

	rtos.IRQ(irq.TIM3).Enable()
}

func blinky(led gpio.Pin, period int) {
	for {
		led.Clear()
		delay.Millisec(100)
		led.Set()
		delay.Millisec(period - 100)
	}
}

func main() {
	//go blinky(led1, 500)
	blinky(led2, 1000)
}

func timISR() {
	t.SR.Store(0)
	led1.Clear()
	delay.Loop(1e3)
	led1.Set()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.TIM3: timISR,
}
