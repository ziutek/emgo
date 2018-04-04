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
	leds  [3]gpio.Pin
	timer *tim.Periph
	ch    = make(chan struct{}, 1)
)

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	gpio.A.EnableClock(false)
	leds[0] = gpio.A.Pin(4)
	leds[1] = gpio.A.Pin(5)
	leds[2] = gpio.A.Pin(9)

	cfg := &gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain}
	for _, led := range leds {
		led.Set()
		led.Setup(cfg)
	}

	timer = tim.TIM3
	pclk := timer.Bus().Clock()
	if pclk < system.AHB.Clock() {
		pclk *= 2
	}
	freq := uint(1e3) // Hz
	timer.EnableClock(true)
	timer.PSC.Store(tim.PSC(pclk/freq - 1))
	timer.ARR.Store(200) // ms
	timer.DIER.Store(tim.UIE)
	timer.CR1.Store(tim.CEN)

	rtos.IRQ(irq.TIM3).Enable()
}

func blinky(led gpio.Pin, period int) {
	for range ch {
		led.Clear()
		delay.Millisec(100)
		led.Set()
		delay.Millisec(period - 100)
	}
}

func main() {
	go blinky(leds[1], 500)
	blinky(leds[2], 500)
}

func timerISR() {
	timer.SR.Store(0)
	leds[0].Set()
	select {
	case ch <- struct{}{}:
		// Success
	default:
		leds[0].Clear()
	}
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.TIM3: timerISR,
}
