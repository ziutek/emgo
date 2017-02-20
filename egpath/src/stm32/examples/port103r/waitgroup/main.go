package main

import (
	"delay"
	"math/rand"
	"rtos"
	"sync"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
)

var (
	leds [4]gpio.Pin
	key3 gpio.Pin
	rnd  rand.XorShift64
)

func init() {
	system.Setup(8, 1, 72/8)
	rtcst.Setup(32768)

	// GPIO

	gpio.B.EnableClock(false)
	leds[2] = gpio.B.Pin(5)
	leds[1] = gpio.B.Pin(6)
	leds[0] = gpio.B.Pin(7)

	gpio.C.EnableClock(true)
	key3 = gpio.C.Pin(10)

	gpio.D.EnableClock(false)
	leds[3] = gpio.D.Pin(2)

	// LEDs

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	for _, pin := range leds {
		pin.Setup(&cfg)
	}

	// Key

	key3.Setup(&gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})

	rnd.Seed(rtos.Nanosec())
}

func waitkey(wg *sync.WaitGroup) {
	for key3.Load() != 0 {
		delay.Millisec(5)
	}
	wg.Add(-1)
}

func setled(wg *sync.WaitGroup, led gpio.Pin, v int) {
	// rnd.Uint32 isn't thread safe but don't care.
	delay.Millisec(100 + int(rnd.Uint32()&0x3ff))
	led.Store(v)
	delay.Millisec(100 + int(rnd.Uint32()&0x3ff))
	wg.Add(-1)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go waitkey(&wg)
	wg.Wait()
	for _, led := range leds[:3] {
		wg.Add(1)
		go setled(&wg, led, 1)
	}
	wg.Wait()
	for _, led := range leds[:3] {
		wg.Add(1)
		go setled(&wg, led, 0)
	}
	wg.Wait()
	leds[3].Set()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,
}
